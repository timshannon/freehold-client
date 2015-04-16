// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import "time"

// Token is the client side defintion to hold the Token
// data returned from a freehold instance.
// a Token is a unique identifier which can grant access to
// a specific resource, or to everything a user has
type Token struct {
	Token      string `json:"token,omitempty"` // Unique identifier used for authentication
	ID         string `json:"id,omitempty"`    //Unique identifier used for identification
	Name       string `json:"name,omitempty"`
	Expires    string `json:"expires,omitempty"`
	Resource   string `json:"resource,omitempty"`
	Permission string `json:"permission,omitempty"`
	Created    string `json:"created,omitempty"`

	client      *Client
	expiresTime time.Time
	createdTime time.Time
}

// AllTokens retrieves all tokens for the given user
// who made the client connection
func (c *Client) AllTokens() ([]*Token, error) {
	var t []*Token
	err := c.doRequest("GET", "/v1/auth/token/", nil, &t)
	if err != nil {
		return nil, err
	}

	for i := range t {
		t[i].client = c
	}

	return t, nil
}

// GetToken retrieves a specific token identified by the passed in id
func (c *Client) GetToken(id string) (*Token, error) {
	t := &Token{}

	err := c.doRequest("GET", "/v1/auth/token/",
		map[string]string{
			"id": id,
		}, &t)

	if err != nil {
		return nil, err
	}
	t.client = c
	return t, nil
}

// NewToken generates a new token with the passed in values.  Only Name is required.
// Depending on the freehold settings this call may need to be made from a client
// which has a users password specified instead of another token
func (c *Client) NewToken(name, resource, permission string, expires time.Time) (*Token, error) {
	t := &Token{
		Name:       name,
		Resource:   resource,
		Permission: permission,
	}

	if !expires.IsZero() {
		t.Expires = expires.Format(time.RFC3339)
	}

	err := c.doRequest("POST", "/v1/auth/token/", t, &t)

	if err != nil {
		return nil, err
	}
	t.client = c
	return t, nil
}

// Delete deletes the current token from the freehold instance
// making it invalid for all future uses
func (t *Token) Delete() error {
	return t.client.doRequest("DELETE", "/v1/auth/token/",
		map[string]string{
			"id": t.ID,
		}, nil)
}

// ExpiresTime is the parsed Time from the Token's JSON string response
func (t *Token) ExpiresTime() time.Time {
	if t.expiresTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, t.Expires)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad Expired date!")
		}
		t.expiresTime = tme
	}
	return t.expiresTime
}

// CreatedTime is the parsed Time from the Token's JSON string response
func (t *Token) CreatedTime() time.Time {
	if t.createdTime.IsZero() {
		tme, err := time.Parse(time.RFC3339, t.Created)
		if err != nil {
			//shouldn't happen as it means freehold is
			// sending out bad dates
			panic("Freehold instance has bad date!")
		}
		t.createdTime = tme
	}
	return t.createdTime
}
