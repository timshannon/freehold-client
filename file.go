// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
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
	Property
}

// GetFile retrieves a file for reading or writing from a freehold instance
func (c *Client) GetFile(filePath string) (*File, error) {
	filePath = strings.TrimSuffix(filePath, "/")
	propPath := propertyPath(filePath)

	f := &File{Property{}}
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

	name := path.Base(info.Name())

	return c.UploadFromReader(name, file, info.Size(), info.ModTime(), dest)

}

// UploadFromReader uploads file data from the passed in reader
// size is required and dest must be a directory on the freehold instance
func (c *Client) UploadFromReader(fileName string, r io.Reader, size int64, modTime time.Time, dest *File) (*File, error) {
	if !dest.IsDir {
		return nil, errors.New("Destination is not a directory.")
	}
	f := &File{
		Property: Property{
			Name:   fileName,
			URL:    path.Join(dest.URL, fileName),
			client: c,
		},
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

// Move moves a file to a new location
func (f *File) Move(to string) error {
	if !strings.HasPrefix(to, "/v1/file/") {
		return errors.New("Invalid file path")
	}
	return f.client.doRequest("PUT", f.URL, map[string]string{"move": to}, nil)
}

// Children returns the child files (if any) of the given folder
// Calling Children on a non-dir file will not error but return
// an empty slice
func (f *File) Children() ([]*File, error) {
	children, err := f.Property.Children()
	if err != nil {
		return nil, err
	}

	files := make([]*File, len(children))
	for i := range children {
		files[i] = &File{children[i]}
	}
	return files, err
}
