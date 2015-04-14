// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Datastore is a datastore stored on freehold instance
// and the properties associated with it
type Datastore struct {
	Property
}

// GetDatastore retrieves a datastore for reading or writing from a freehold instance
func (c *Client) GetDatastore(filePath string) (*Datastore, error) {
	filePath = strings.TrimSuffix(filePath, "/")
	propPath := propertyPath(filePath)

	d := &Datastore{Property{}}
	err := c.doRequest("GET", propPath, nil, d)

	if err != nil {
		return nil, err
	}

	d.client = c

	return d, nil
}

// NewDatastore creates a new datastore file at the path, passed in
func (c *Client) NewDatastore(filePath string) (*Datastore, error) {
	err := c.doRequest("POST", filePath, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.GetDatastore(filePath)
}

// UploadDatastore uploads a local datstore file to the freehold instance
// and returns a Datastore
// Dest must be a Dir
func (c *Client) UploadDatastore(dsFile *os.File, dest *File) (*Datastore, error) {
	info, err := dsFile.Stat()
	if err != nil {
		return nil, err
	}

	name := filepath.Base(info.Name())
	if !dest.IsDir {
		return nil, errors.New("Destination is not a directory.")
	}

	d := &Datastore{
		Property: Property{
			Name:   name,
			URL:    filepath.Join(dest.URL, name),
			client: c,
		},
	}

	err = d.upload("POST", dsFile, info.Size(), info.ModTime())
	if err != nil {
		return nil, err
	}

	return c.GetDatastore(d.URL)
}

// Children returns the child datastores (if any) of the given folder
// Calling Children on a non-dir file will not error but return
// an empty slice
func (d *Datastore) Children() ([]*File, error) {
	children, err := d.Property.Children()
	if err != nil {
		return nil, err
	}

	files := make([]*File, len(children))
	for i := range children {
		files[i] = &File{children[i]}
	}
	return files, err
}

// Get gets a value out of a freehold datastore
func (d *Datastore) Get(key, returnValue interface{}) error {
	return d.client.doRequest("GET", d.URL, map[string]interface{}{
		"key": key,
	}, returnValue)
}

// Min gets the minimum key / value in the datastore
//func (d *Datastore) Min(key, value interface{}) error {
//return d.client.doRequest("GET", d.URL, map[string}struct{}
//"min": struct{}{},
//}, returnValue)
//}
