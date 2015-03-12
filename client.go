// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is used for interacting with a Freehold Instance
// A Client needs a url, username, and password or token
// After the client is initialized, all requests should be run
// against the path of the file, ds, etc in question
// /v1/file/test.txt instead of the full url
// https://freeholdinstance/v1/file/test.txt
type Client struct {
	hClient  *http.Client
	root     *url.URL
	username string
	pass     string
}

//jsend is the reponse format from a freehold instance
type jsend struct {
	Status   string           `json:"status"`
	Data     *json.RawMessage `json:"data,omitempty"`
	Message  string           `json:"message,omitempty"`
	Failures []error          `json:"failures,omitempty"`
}

// New creates a new Freehold Client
// tlsCfg is optional
func New(rootURL, username, passwordOrToken string, tlsCfg *tls.Config) (*Client, error) {
	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
	}

	uri, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		root:     uri,
		username: username,
		pass:     passwordOrToken,
		hClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
		},
	}
	return c, nil
}

//doRequest will run a standard freehold request, and try to unpack the data result into
// the passed in result interface.  Only to be used with
// JSEND responses
func (c *Client) doRequest(method, fhPath string, result interface{}) error {
	c.root.Path = fhPath
	req, err := http.NewRequest(method, c.root.String(), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.pass)

	res, err := c.hClient.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	response := &jsend{}
	err = decoder.Decode(response)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Request %s failed with a status of %d.  Message: %s", c.root.String(), res.StatusCode, response.Message)
	}

	err = json.Unmarshal(*response.Data, result)
	if err != nil {
		return err
	}
	return nil
}
