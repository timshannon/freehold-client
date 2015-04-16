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

func TestApps(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/application/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				body := requestBody(t, r)
				if strings.Contains(body, "id") {
					//one user
					fmt.Fprint(w, `{"status":"success","data":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"}}`)
				} else {
					// all users
					fmt.Fprint(w, `{"status":"success","data":{"admin":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"},"datastore":{"id":"datastore","name":"Datastore Viewer","description":"Datastore Viewer - Views datastore files","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/v1/file/index.html","icon":"/datastore/v1/file/image/datastore_icon.png","version":"0.1","file":"datastore.zip"}}}`)
				}
			}

			if r.Method == "POST" {
				// New user
				fmt.Fprint(w, `{"status":"success","data":{"name":"test","homeApp":"home"}}`)
			}
			if r.Method == "PUT" {
				// update user
				fmt.Fprint(w, `{"status":"success"}`)
			}
			if r.Method == "DELETE" {
				// delete user
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	client, err := New(server.URL, username, password)
	//client, err := New("https://tshannon.org", "tshannon", "_hvlkuhuCsGifYbxVBuEgaLbYvAb9kVz6dHA43XThCk=")
	if err != nil {
		t.Fatal(err)
	}

	all, err := client.AllApplications()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for k := range all {
		if all[k].ID == "admin" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Admin app not found!")
	}

	_, err = client.GetApplication("admin")
	if err != nil {
		t.Fatal(err)
	}

}
