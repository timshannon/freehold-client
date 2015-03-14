// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package freeholdclient is a Go client for interacting with a freehold instance.
//When using the freehold client it is recommended to generate a token and use that
//rather than storing a users password locally in cleartext
package freeholdclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is used for interacting with a Freehold Instance
// A Client needs a url, username, and password or token
// After the client is initialized, all requests should be run
// against the path of the file, ds, etc in question
// /v1/file/test.txt instead of the full url
// https://freeholdinstance/v1/file/test.txt
//
// It is highly encouraged to not store a user's password, and
// instead use a Security Token generated for this specific
// client
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
func (c *Client) doRequest(method, fhPath string, send interface{}, result interface{}) error {
	c.root.Path = fhPath

	req, err := http.NewRequest(method, c.root.String(), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.pass)
	if send != nil {
		b, err := json.Marshal(send)
		if err != nil {
			return fmt.Errorf("Error json marshalling send data: %v", err)
		}
		r := bytes.NewReader(b)
		req.Body = ioutil.NopCloser(r)
		req.ContentLength = int64(r.Len())
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

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

	err = isError(c.root.String(), res.StatusCode, response)
	if err != nil {
		return err
	}

	if result != nil {
		err = json.Unmarshal(*response.Data, result)
		if err != nil {
			return err
		}
	}
	return nil
}
