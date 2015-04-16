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

func TestUser(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/user/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				body := requestBody(t, r)
				if strings.Contains(body, "user") {
					//one user
					fmt.Fprint(w, `{"status":"success","data":{"name":"not an admin","homeApp":"home"}}`)
				} else {
					// all users
					fmt.Fprint(w, `{"status":"success","data":{"tester":{"name":"not an admin","homeApp":"home"},
						"quinitTestUser":{"name":"Bob QUnit Test User","homeApp":"home"},
						"tshannon":{"name":"Tim Shannon","homeApp":"home","admin":true},
						"tshannon2":{"name":"Tester","homeApp":"home","admin":true}}}`)
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
	if err != nil {
		t.Fatal(err)
	}

	all, err := client.AllUsers()
	if err != nil {
		t.Fatal(err)
	}

	if len(all) != 4 {
		t.Errorf("User count doesn't match expected 4, got %d", len(all))
	}

	nu, err := client.NewUser("test", "testtesttest", "test", "home", false)
	if err != nil {
		t.Fatal(err)
	}
	if nu.Name != "test" {
		t.Errorf("Username doesn't match, expected %s, got %s", "test", nu.Name)
	}

	u, err := client.GetUser("test")
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "test" {
		t.Errorf("Username doesn't match, expected %s, got %s", "test", u.Username)
	}

	err = u.SetName("Test New Name")
	if err != nil {
		t.Fatal(err)
	}

	err = u.Delete()
	if err != nil {
		t.Fatal(err)
	}
}
