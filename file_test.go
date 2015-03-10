package freeholdclient

import (
	"crypto/tls"
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
	filePath = "/v1/file/testing/"
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
			fmt.Fprint(w, `{"status":"success","data":{"name":"testing","url":"/v1/file/testing/","permissions":{"owner":"tshannon","private":"rw"},"modified":"2015-03-06T15:47:40-06:00","isDir":true}}`)
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	f, err := client.RetrieveFile(filePath)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(f.Permissions.Owner)

}

func TestRealFile(t *testing.T) {
	client, err := New("https://tshannon.org", "tshannon", "eAhmtD1hJfzV4C3GvDSIOAkumN54vSgCOIOcrLb5w2A=",
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		t.Fatal(err)
	}

	f, err := client.RetrieveFile(filePath)

	if err != nil {
		t.Fatal(err)
	}
	children, err := f.Children()
	if err != nil {
		t.Fatal(err)
	}
	for i := range children {
		t.Log(children[i].URL)
	}
}
