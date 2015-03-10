// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// File is a file stored on freehold instance
// and the properties associated with it
type File struct {
	Name        string      `json:"name,omitempty"`
	URL         string      `json:"url,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
	Size        int64       `json:"size,omitempty"`
	Modified    string      `json:"modified,omitempty"`
	IsDir       bool        `json:"isDir,omitempty"`
	client      *Client
	mpReader    *multipart.Reader
}

// Permission is the client side definition of a Freehold Permission
type Permission struct {
	Owner   string `json:"owner,omitempty"`
	Public  string `json:"public,omitempty"`
	Friend  string `json:"friend,omitempty"`
	Private string `json:"private,omitempty"`
}

// RetrieveFile retrieves a file from a freehold instance
func (c *Client) RetrieveFile(filePath string) (*File, error) {
	filePath = strings.TrimSuffix(filePath, "/")
	propPath := propertyPath(filePath)

	f := &File{}
	err := c.doRequest("GET", propPath, f)

	if err != nil {
		return nil, err
	}

	f.client = c

	return f, nil
}

// ModifiedDate is the parsed date time from the modified string value
// returned from the freehold REST API
func (f *File) ModifiedDate() time.Time {
	t, err := time.Parse(time.RFC3339, f.Modified)
	if err != nil {
		//shouldn't happen as it means freehold is
		// sending out bad dates
		panic("Freehold instance is has bad Modified date!")
	}
	return t
}

// Children returns the child files (if any) of the given file
// Calling Children on a non-dir file will not error but return
// an empty slice
func (f *File) Children() ([]*File, error) {
	if !f.IsDir {
		return []*File{}, nil
	}

	uri := propertyPath(f.URL)
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}

	var children []*File
	err := f.client.doRequest("GET", uri, &children)
	if err != nil {
		return nil, err
	}
	return children, nil
}

// Reads data from the freehold instance on the given file (GET file data)
func (f *File) Read(p []byte) (n int, err error) {
	if f.mpReader == nil {
		req, err := http.NewRequest(method, f.URL, nil)
		if err != nil {
			return 0, err
		}

		req.SetBasicAuth(f.client.username, f.client.pass)

		res, err := f.client.Do(req)
		defer res.Body.Close() //TODO: How to handle this with mpReader?
		if err != nil {
			return 0, err
		}

		mediaType, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
		if err != nil {
			return 0, err
		}

		if !strings.HasPrefix(mediaType, "multipart/") {
			return 0, errors.New("Invalid multipart form response from ", f.URL)
		}
		f.mpReader = multipart.NewReader(res.Body, params["boundary"])
	}

	part, err := f.mpReader.NextPart()
	if err == io.EOF {
		return 0, nil
	}
	if err != nil {
		return f.mpReaderLen, err
	}

	buff := bytes.NewBuffer(p)
	n, err := io.Copy(buff, part)
	if err != nil {
		return n, err
	}
	return n, part.Close()
}

// Close closes the open file
func (f *File) Close() error {
}

// Write writes data to the freehold file (updates it with a Put Call)
func (f *File) Write(b []byte) (n int, err error) {

}
