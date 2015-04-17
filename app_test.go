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
					//one app
					fmt.Fprint(w, `{"status":"success","data":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"}}`)
				} else {
					// all apps
					fmt.Fprint(w, `{"status":"success","data":{"admin":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"},"datastore":{"id":"datastore","name":"Datastore Viewer","description":"Datastore Viewer - Views datastore files","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/v1/file/index.html","icon":"/datastore/v1/file/image/datastore_icon.png","version":"0.1","file":"datastore.zip"}}}`)
				}
			}

			if r.Method == "PUT" || r.Method == "POST" {
				// install / Upgrade App
				fmt.Fprint(w, `{"status":"success","data":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"}}`)
			}
			if r.Method == "DELETE" {
				// delete app
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	// Available
	mux.HandleFunc("/v1/application/available/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				// all apps
				fmt.Fprint(w, `{"status":"success","data":{"admin":{"id":"admin","name":"Admin Console","description":"Administrator's Console. For managing users, logs, and freehold settings.","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/admin/v1/file/index.html","icon":"/admin/v1/file/image/admin_icon.png","version":"0.1","file":"admin.zip"},"datastore":{"id":"datastore","name":"Datastore Viewer","description":"Datastore Viewer - Views datastore files","author":"Tim Shannon - shannon.timothy@gmail.com","root":"/v1/file/index.html","icon":"/datastore/v1/file/image/datastore_icon.png","version":"0.1","file":"datastore.zip"}}}`)
			}

			if r.Method == "POST" {
				// post new available app
				fmt.Fprint(w, `{"status":"success","data":"blog.zip"}`)
			}
		})

	client, err := New(server.URL, username, password)
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

	a, err := client.GetApplication("admin")
	if err != nil {
		t.Fatal(err)
	}

	err = a.Uninstall()
	if err != nil {
		t.Fatal(err)
	}

	available, err := client.AvailableApplications()
	if err != nil {
		t.Fatal(err)
	}

	var avInstall *AvailableApplication
	found = false
	for k := range available {
		if available[k].ID == "admin" {
			avInstall = available[k]
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Admin app not found in available apps!")
	}

	a, err = avInstall.Install()
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != "admin" {
		t.Fatalf("Admin app not installed!")
	}

	a, err = avInstall.Upgrade()
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != "admin" {
		t.Fatalf("Admin app not upgraded!")
	}

	avInstall, err = client.PostAvailableApplication("http://test.com/app.zip")
	if err != nil {
		t.Fatal(err)
	}

}
