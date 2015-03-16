// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"bytes"
	"errors"
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
	modTime    time.Time
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
	err := c.doRequest("GET", propPath, nil, f)

	if err != nil {
		return nil, err
	}

	f.client = c

	return f, nil
}

// NewFolder creates a new folder on the freehold instance
func (c *Client) NewFolder(folderPath string) error {
	if !strings.HasPrefix(folderPath, "/v1/file/") {
		return errors.New("Invalid folder path")
	}
	return c.doRequest("POST", folderPath, nil, nil)
}

// UploadFile uploads a local file to the freehold instance
// and returns a File type.
// Dest must be a Dir
func (c *Client) UploadFile(file *os.File, dest *File) (*File, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	name := filepath.Base(info.Name())

	return c.UploadFromReader(name, file, info.Size(), info.ModTime(), dest)

}

// UploadFromReader uploads file data from the passed in reader
// size is required and dest must be a directory on the freehold instance
func (c *Client) UploadFromReader(fileName string, r io.Reader, size int64, modTime time.Time, dest *File) (*File, error) {
	if !dest.IsDir {
		return nil, errors.New("Destination is not a directory.")
	}
	f := &File{
		Name:   fileName,
		URL:    filepath.Join(dest.URL, fileName),
		client: c,
	}

	err := f.upload("POST", r, size, modTime)
	if err != nil {
		return nil, err
	}

	return c.GetFile(f.URL)
}

// Update overwrites the given file with the bytes read from r
// Size is the total size to be read from r, and a limitReader is used to
// enforce this
func (f *File) Update(r io.Reader, size int64) error {
	return f.upload("PUT", r, size, time.Time{})
}

func (f *File) upload(method string, r io.Reader, size int64, modTime time.Time) error {
	lr := io.LimitReader(r, size)

	var res *http.Response

	pRead, pWrite := io.Pipe()
	writer := multipart.NewWriter(pWrite)

	done := make(chan error, 1)

	f.client.root.Path = path.Dir(f.URL)
	uri := f.client.root.String()

	go func() {
		defer pWrite.Close()
		prt, err := writer.CreateFormFile("file", f.Name)
		if err != nil {
			done <- err
			return
		}

		_, err = io.Copy(prt, lr)
		if err == nil {
			err = writer.Close()
		}

		done <- err
	}()

	req, err := http.NewRequest(method, uri, pRead)
	if err != nil {
		return err
	}

	req.SetBasicAuth(f.client.username, f.client.pass)
	if !modTime.IsZero() {
		req.Header.Set("Fh-Modified", modTime.Format(time.RFC3339))
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = multipartOverhead + size + int64(len([]byte("file"+f.Name)))

	res, err = f.client.hClient.Do(req)

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
func (f *File) ModifiedTime() time.Time {
	if f.modTime.IsZero() {
		t, err := time.Parse(time.RFC3339, f.Modified)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance is has bad Modified date!")
		}
		f.modTime = t
	}
	return f.modTime
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
	err := f.client.doRequest("GET", uri, nil, &children)
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

		err = isError(f.FullURL(), res.StatusCode, nil)
		if err != nil {
			return 0, err
		}

		if res != nil {
			f.readerBody = res.Body
		}

	}
	return f.readerBody.Read(p)
}

// Close closes the open reader
func (f *File) Close() error {
	if f.readerBody != nil {
		r := f.readerBody
		f.readerBody = nil
		return r.Close()
	}
	return nil
}

// Delete deletes a file
func (f *File) Delete() error {
	return f.client.doRequest("DELETE", f.URL, nil, nil)
}

// Move moves a file to a new location
func (f *File) Move(to string) error {
	if !strings.HasPrefix(to, "/v1/file/") {
		return errors.New("Invalid file path")
	}
	return f.client.doRequest("PUT", f.URL, map[string]string{"move": to}, nil)
}
