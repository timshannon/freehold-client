// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// multipartOverhead is how many extra bytes mime/multipart's
// Writer adds around content
// Thanks camlistore - https://github.com/camlistore/camlistore/blob/master/pkg/client/upload.go
var multipartOverhead = func() int64 {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, _ := w.CreateFormFile("0", "0")

	dummyContents := []byte("0")
	part.Write(dummyContents)

	w.Close()
	return int64(b.Len()) - 3 // remove what was added
}()

// File is a file stored on freehold instance
// and the properties associated with it
type File struct {
	Name        string      `json:"name,omitempty"`
	URL         string      `json:"url,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
	Size        int64       `json:"size,omitempty"`
	Modified    string      `json:"modified,omitempty"`
	IsDir       bool        `json:"isDir,omitempty"`

	client     *Client
	readerBody io.ReadCloser
}

// Permission is the client side definition of a Freehold Permission
type Permission struct {
	Owner   string `json:"owner,omitempty"`
	Public  string `json:"public,omitempty"`
	Friend  string `json:"friend,omitempty"`
	Private string `json:"private,omitempty"`
}

// GetFile retrieves a file for reading or writing from a freehold instance
func (c *Client) GetFile(filePath string) (*File, error) {
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

// UploadFile uploads a local file to the freehold instance
// and returns a File type.
// Dest must be a Dir
func (c *Client) UploadFile(file *os.File, dest *File) (*File, error) {
	if !dest.IsDir {
		return nil, errors.New("Destination is not a directory.")
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	name := filepath.Base(info.Name())
	f := &File{
		Name:   name,
		URL:    filepath.Join(dest.URL, name),
		client: c,
	}

	err = f.upload("POST", file, info.Size())
	if err != nil {
		return nil, err
	}

	return c.GetFile(f.URL)
}

// Update overwrites the given file with the bytes read from r
// Size is the total size to be read from r, and a limitReader is used to
// enforce this
func (f *File) Update(r io.Reader, size int64) error {
	return f.upload("PUT", r, size)
}

func (f *File) upload(method string, r io.Reader, size int64) error {
	lr := io.LimitReader(r, size)

	var res *http.Response

	pout, pin := io.Pipe()
	writer := multipart.NewWriter(pin)
	defer writer.Close()
	done := make(chan error)

	uri := path.Dir(f.FullURL())

	go func() {
		req, err := http.NewRequest(method, uri, pout)
		if err != nil {
			done <- err
			return
		}

		req.SetBasicAuth(f.client.username, f.client.pass)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ContentLength = multipartOverhead + size

		res, err = f.client.hClient.Do(req)

		if err != nil {
			done <- err
			return
		}

		if res.StatusCode >= 400 {
			done <- fmt.Errorf("Request %s failed with a status of %d.", uri, res.StatusCode)
			return
		}
		done <- nil
	}()

	fmt.Println("before createform file")
	prt, err := writer.CreateFormFile("file", f.Name)

	defer pin.Close()
	if err != nil {
		return err
	}

	fmt.Println("after create form file")
	_, err = io.Copy(prt, lr)
	fmt.Println("after copy before error")
	if err != nil {
		return err
	}

	fmt.Println("after copy after error")

	return <-done
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

// FullURL returns the full url of the file including the
// root of the freehold instance
func (f *File) FullURL() string {
	f.client.root.Path = f.URL
	return f.client.root.String()
}

// Reads data from the freehold instance on the given file (GET file data)
// Close() needs to be called when read is completed
func (f *File) Read(p []byte) (n int, err error) {
	if f.readerBody == nil {
		req, err := http.NewRequest("GET", f.FullURL(), nil)
		if err != nil {
			return 0, err
		}

		req.SetBasicAuth(f.client.username, f.client.pass)
		res, err := f.client.hClient.Do(req)

		if err != nil {
			return 0, err
		}

		if res.StatusCode != 200 {
			return 0, fmt.Errorf("File Retrieve %s failed with a status of %d.", f.URL, res.StatusCode)
		}

		if res != nil {
			f.readerBody = res.Body
		}

	}
	return f.readerBody.Read(p)
}

// Close closes the open file
func (f *File) Close() error {
	if f.readerBody != nil {
		r := f.readerBody
		f.readerBody = nil
		return r.Close()
	}
	return nil
}

//TODO: https://gist.github.com/cryptix/9dd094008b6236f4fc57
