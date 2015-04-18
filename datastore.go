// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
)

// Datastore is a datastore stored on freehold instance
// and the properties associated with it
type Datastore struct {
	Property
}

// KeyValue is a key value pair returned from a datastore
type KeyValue struct {
	K           *json.RawMessage `json:"key,omitempty"`
	V           *json.RawMessage `json:"value,omitempty"`
	errRetrieve error            // any error returned when retrieving this KV
}

// Key is a convience function for unmarshalling the key
// from the KeyValue, will return any errors from the retrieval
// of this key first
// So you can run requests like err := ds.Min().Key(&result)
// and consolidate your error checking into one call
func (kv *KeyValue) Key(result interface{}) error {
	if kv.errRetrieve != nil {
		return kv.errRetrieve
	}
	return json.Unmarshal([]byte(*kv.K), result)
}

// Value is a convience function for unmarshalling the key
// from the KeyValue, will return any errors from the retrieval
// of this key first
// So you can run requests like err := ds.Min().Value(&result)
// and consolidate your error checking into one call
func (kv *KeyValue) Value(result interface{}) error {
	if kv.errRetrieve != nil {
		return kv.errRetrieve
	}
	return json.Unmarshal([]byte(*kv.V), result)
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

	name := path.Base(info.Name())
	if !dest.IsDir {
		return nil, errors.New("Destination is not a directory.")
	}

	d := &Datastore{
		Property: Property{
			Name:   name,
			URL:    path.Join(dest.URL, name),
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

// Drop deletes the datastore file
func (d *Datastore) Drop() error {
	return d.Property.Delete()
}

// Get gets a value out of a freehold datastore
func (d *Datastore) Get(key, returnValue interface{}) error {
	return d.client.doRequest("GET", d.URL, map[string]interface{}{
		"key": key,
	}, returnValue)
}

// Put puts a new key value pair into the datastore
func (d *Datastore) Put(key, value interface{}) error {
	return d.client.doRequest("PUT", d.URL, map[string]interface{}{
		"key":   key,
		"value": value,
	}, nil)
}

// PutObj puts the entire passed in object into the datastore
// Top level keys become the basis for the key / values
// object must be able to be marshalled into a json string
func (d *Datastore) PutObj(object interface{}) error {
	return d.client.doRequest("PUT", d.URL, object, nil)

}

// Delete deletes the value from the datastore for the passed in key
func (d *Datastore) Delete(key interface{}) error {
	return d.client.doRequest("DELETE", d.URL, map[string]interface{}{
		"key": key,
	}, nil)
}

// Min gets the minimum key / value in the datastore
// Example:
// 	err := ds.Min().Value(&result)
func (d *Datastore) Min() *KeyValue {
	result := &KeyValue{}
	err := d.client.doRequest("GET", d.URL, map[string]struct{}{
		"min": struct{}{},
	}, result)
	if err != nil {
		result.errRetrieve = err
	}
	return result
}

// Max gets the maximum key / value in the datastore
// Example:
// 	err := ds.Max().Key(&result)
func (d *Datastore) Max() *KeyValue {
	result := &KeyValue{}
	err := d.client.doRequest("GET", d.URL, map[string]struct{}{
		"max": struct{}{},
	}, result)
	if err != nil {
		result.errRetrieve = err
	}
	return result
}

// Iter is for iterating through a datastore
type Iter struct {
	From   interface{} `json:"from,omitempty"`
	To     interface{} `json:"to,omitempty"`
	Skip   int         `json:"skip,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Regexp string      `json:"regexp,omitempty"`
	Order  string      `json:"order,omitempty"`
}

// Iter returns the list of key / values matched by the passed in interator
func (d *Datastore) Iter(iter *Iter) ([]*KeyValue, error) {
	var result []*KeyValue
	err := d.client.doRequest("GET", d.URL, map[string]interface{}{
		"iter": iter,
	}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
