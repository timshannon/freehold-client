// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// File is a file stored on freehold instance
// and the properties associated with it
type File struct {
	Name        string      `json:"name,omitempty"`
	URL         string      `json:"url,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
	Size        int64       `json:"size,omitempty"`
	Modified    string      `json:"modified,omitempty"`
	IsDir       bool        `json:"isDir,omitempty"`
}

type Permission struct {
	Owner   string `json:"owner,omitempty"`
	Public  string `json:"public,omitempty"`
	Friend  string `json:"friend,omitempty"`
	Private string `json:"private,omitempty"`
}

// RetrieveFile retrieves a file from a freehold instance
func (c *Client) RetrieveFile(path string) (*File, error) {

}
