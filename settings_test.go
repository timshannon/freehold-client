// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestSettings(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/settings/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				body := requestBody(t, r)
				if strings.Contains(body, "404File") {
					fmt.Fprint(w, `{"status":"success","data":{"404File":{"description":"Path to a standard 404 page.","value":"/core/v1/file/404.html"}}}`)
				} else {
					fmt.Fprint(w, `{"status":"success","data":{"404File":{"description":"Path to a standard 404 page.","value":"/core/v1/file/404.html"},"AllowWebAppInstall":{"description":"Whether or not applications are allowed to be installed from any arbitary url.  i.e http://github.com/developer/app/app.zip","value":true},"DatastoreFileTimeout":{"description":"The number of seconds of inactivity needed before a datastore file is automatically closed. The higher the timeout the more resources needed to hold open more files.  The lower the timeout the more clients will be waiting on locked files to be opened.","value":60}}}`)
				}
			}

			if r.Method == "DELETE" {
				// default setting
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	client, err := New(server.URL, username, password)
	if err != nil {
		t.Fatal(err)
	}

	s, err := client.AllSettings()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for k := range s {
		if k == "404File" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Setting 404File not found!")
	}

	_, err = client.GetSetting("404File")
	if err != nil {
		t.Fatal(err)
	}

	if client.DefaultSetting("404File") != nil {
		t.Fatal(err)
	}

}
