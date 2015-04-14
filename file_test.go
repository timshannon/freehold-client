// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	username = "tester"
	password = "testerToken"
	dirPath  = "/v1/file/testing/"
)

func startMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func stopMockServer() {
	server.Close()
}

func TestFile(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/properties/file/testing",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":{"name":"testing","url":"/v1/file/testing/",
				"permissions":{"owner":"tshannon","private":"rw"},"modified":"2015-03-06T15:47:40-06:00","isDir":true}}`)
		})
	mux.HandleFunc("/v1/properties/file/testing/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":[{"name":"test.txt","url":"/v1/file/testing/test.txt",
				"permissions":{"owner":"tshannon","private":"rw"},"size":9,"modified":"2015-03-13T11:28:59-05:00"}]}`)
		})

	client, err := New(server.URL, username, password)
	if err != nil {
		t.Fatal(err)
	}

	f, err := client.GetFile(dirPath)

	if err != nil {
		t.Fatalf("Error retrieving file properties: %v", err)
	}

	if f.Permissions.Owner != "tshannon" {
		t.Errorf("Permissions owner does not match.  Expected tshannon got %s", f.Permissions.Owner)
	}

	if f.Permissions.Private != "rw" {
		t.Errorf("Permissions Private does not match.  Expected rw got %s", f.Permissions.Private)
	}

	if m, _ := time.Parse(time.RFC3339, "2015-03-06T15:47:40-06:00"); !f.ModifiedTime().Equal(m) {
		t.Errorf("Modified Date does not match. Expected %v got %v", m, f.ModifiedTime())
	}

	children, err := f.Children()
	if err != nil {
		t.Fatalf("Error retrieving child files: %v", err)
	}
	if children[0].Name != "test.txt" {
		t.Errorf("Child file name does not match. Expected test.txt got %s", children[0].Name)
	}

}
