// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	username = "tester"
	password = "testerToken"
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
	mux.HandleFunc("/v1/auth/token/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":{"name":"testing","url":"/v1/file/testing/",
				"permissions":{"owner":"tshannon","private":"rw"},"modified":"2015-03-06T15:47:40-06:00","isDir":true}}`)
		})
	mux.HandleFunc("/v1/properties/file/testing/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":[{"name":"test.txt","url":"/v1/file/testing/test.txt",
				"permissions":{"owner":"tshannon","private":"rw"},"size":9,"modified":"2015-03-13T11:28:59-05:00"}]}`)
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	f, err := client.GetFile(dirPath)

}
