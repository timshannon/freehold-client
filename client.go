// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"crypto/tls"
	"net/http"
)

// Client is used for interacting with a Freehold Instance
// A Client needs a url, username, and password or token
type Client struct {
	*http.Client
	rootURL  string
	username string
	password string
	token    string
}

// New creates a new Freehold Client
// tlsCfg is optional
func New(rootURL, username, password, token string, tlsCfg *tls.Config) *Client {
	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
	}
	return &Client{
		rootURL:  rootURL,
		username: username,
		token:    token,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
		},
	}
}
