// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

// Auth contains the type and identity of a user in Freehold
// if user == nil, then auth is public access
type Auth struct {
	AuthType string `json:"type"`
	Username string `json:"user,omitempty"`
	*User
	*Token
}

// Auth returns authention information about the current user
func (c *Client) Auth() (*Auth, error) {
	a := &Auth{}
	err := c.doRequest("GET", "/v1/auth/", nil, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
