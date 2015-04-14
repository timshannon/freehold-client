// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"
)

// Property is a  set of datastore or file properties
// from a freehold instance
type Property struct {
	Name        string      `json:"name,omitempty"`
	URL         string      `json:"url,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
	Size        int64       `json:"size,omitempty"`
	Modified    string      `json:"modified,omitempty"`
	IsDir       bool        `json:"isDir,omitempty"`

	client     *Client
	readerBody io.ReadCloser
	modTime    time.Time
}

// Permission is the client side definition of a Freehold Permission
type Permission struct {
	Owner   string `json:"owner,omitempty"`
	Public  string `json:"public,omitempty"`
	Friend  string `json:"friend,omitempty"`
	Private string `json:"private,omitempty"`
}

func (p *Property) upload(method string, r io.Reader, size int64, modTime time.Time) error {
	lr := io.LimitReader(r, size)

	var res *http.Response

	pRead, pWrite := io.Pipe()
	writer := multipart.NewWriter(pWrite)

	done := make(chan error, 1)

	p.client.root.Path = path.Dir(p.URL)
	uri := p.client.root.String()

	go func() {
		defer pWrite.Close()
		prt, err := writer.CreateFormFile("file", p.Name)
		if err != nil {
			done <- err
			return
		}

		written, err := io.Copy(prt, lr)
		if err == nil {
			err = writer.Close()
		}

		if err == nil && written != size {
			err = io.ErrShortWrite
		}

		done <- err
	}()

	req, err := http.NewRequest(method, uri, pRead)
	if err != nil {
		return err
	}

	req.SetBasicAuth(p.client.username, p.client.pass)
	if !modTime.IsZero() {
		req.Header.Set("Fh-Modified", modTime.Format(time.RFC3339))
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = multipartOverhead + size + int64(len([]byte("file"+p.Name)))

	res, err = p.client.hClient.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = isError(uri, res.StatusCode, nil)
	if err != nil {
		return err
	}

	err = <-done
	return err
}

// ModifiedTime is the parsed date and time from the modified string value
// returned from the freehold REST API
func (p *Property) ModifiedTime() time.Time {
	if p.modTime.IsZero() && p.Modified != "" {
		t, err := time.Parse(time.RFC3339, p.Modified)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has a bad Modified date! " + err.Error())
		}
		p.modTime = t
	}
	return p.modTime
}

// Children returns the child files / datastores (if any) of the given folder
// Calling Children on a non-dir file will not error but return
// an empty slice
func (p *Property) Children() ([]Property, error) {
	if !p.IsDir {
		return []Property{}, nil
	}

	uri := propertyPath(p.URL)
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}

	var children []Property
	err := p.client.doRequest("GET", uri, nil, &children)
	if err != nil {
		return nil, err
	}

	for i := range children {
		children[i].client = p.client
	}
	return children, nil
}

// FullURL returns the full url of the file / datstore including the
// root of the freehold instance
func (p *Property) FullURL() string {
	p.client.root.Path = p.URL
	return p.client.root.String()
}

// Reads data from the freehold instance on the given file or datastore (GET file data)
// Close() needs to be called when read is completed
func (p *Property) Read(b []byte) (n int, err error) {
	if p.readerBody == nil {
		req, err := http.NewRequest("GET", p.FullURL(), nil)
		if err != nil {
			return 0, err
		}

		req.SetBasicAuth(p.client.username, p.client.pass)
		res, err := p.client.hClient.Do(req)

		if err != nil {
			return 0, err
		}

		err = isError(p.FullURL(), res.StatusCode, nil)
		if err != nil {
			return 0, err
		}

		if res != nil {
			p.readerBody = res.Body
		}

	}
	return p.readerBody.Read(b)
}

// Close closes the open reader
func (p *Property) Close() error {
	if p.readerBody != nil {
		r := p.readerBody
		p.readerBody = nil
		return r.Close()
	}
	return nil
}

// Delete deletes a file / datastore
func (p *Property) Delete() error {
	return p.client.doRequest("DELETE", p.URL, nil, nil)
}

// SetPermission sets the current file / datastore's permissions to those
// passed in
func (p *Property) SetPermission(prm *Permission) error {
	return p.client.doRequest("PUT", p.URL, map[string]*Permission{"permissions": prm}, nil)
}
