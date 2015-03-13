// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAllTokens(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/token/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":[
				{"id":"f320dbd99401b016c8514132f500735c","name":"testtoken","expires":"2015-06-11T15:17:29-05:00","created":"2015-03-13T15:17:29-05:00"},
				{"id":"61c78f8bf7f62f2b17aaaa4e1345e27b","name":"test","expires":"2015-06-11T14:51:36-05:00","created":"2015-03-13T14:51:36-05:00"},
				{"id":"bf0ac157dbb199d3b2ab95332df7f037","name":"test2","expires":"2015-06-11T14:52:07-05:00","created":"2015-03-13T14:52:07-05:00"}]}`)
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	tkns, err := client.AllTokens()
	if err != nil {
		t.Fatal(err)
	}
	if len(tkns) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(tkns))
	}

	found := false
	for i := range tkns {
		if tkns[i].Name == "test" {
			found = true
			if c, _ := time.Parse(time.RFC3339, "2015-03-13T14:51:36-05:00"); !tkns[i].CreatedTime().Equal(c) {
				t.Errorf("Created time doesn't match for test token. Expected %v got %v", c, tkns[i].CreatedTime())
			}
			if e, _ := time.Parse(time.RFC3339, "2015-06-11T14:51:36-05:00"); !tkns[i].ExpiresTime().Equal(e) {
				t.Errorf("Expired time doesn't match for test token. Expected %v got %v", e, tkns[i].ExpiresTime())
			}
			if tkns[i].ID != "61c78f8bf7f62f2b17aaaa4e1345e27b" {
				t.Errorf("ID doesn't match, expected 61c78f8bf7f62f2b17aaaa4e1345e27b, got %s", tkns[i].ID)
			}
		}
	}

	if !found {
		t.Fatalf("Token named 'test' not found in token list")
	}

}

func TestGetToken(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/token/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":
				{"id":"61c78f8bf7f62f2b17aaaa4e1345e27b","name":"test","expires":"2015-06-11T14:51:36-05:00","created":"2015-03-13T14:51:36-05:00"}}`)
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	tkn, err := client.GetToken("")
	if err != nil {
		t.Fatal(err)
	}

	if tkn.Name != "test" {
		t.Errorf("Name doesn't match, expected test, got %s", tkn.Name)
	}
	if c, _ := time.Parse(time.RFC3339, "2015-03-13T14:51:36-05:00"); !tkn.CreatedTime().Equal(c) {
		t.Errorf("Created time doesn't match for test token. Expected %v got %v", c, tkn.CreatedTime())
	}
	if e, _ := time.Parse(time.RFC3339, "2015-06-11T14:51:36-05:00"); !tkn.ExpiresTime().Equal(e) {
		t.Errorf("Expired time doesn't match for test token. Expected %v got %v", e, tkn.ExpiresTime())
	}
	if tkn.ID != "61c78f8bf7f62f2b17aaaa4e1345e27b" {
		t.Errorf("ID doesn't match, expected 61c78f8bf7f62f2b17aaaa4e1345e27b, got %s", tkn.ID)
	}

}

func TestNewToken(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/token/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":
				{"id":"61c78f8bf7f62f2b17aaaa4e1345e27b","name":"test","expires":"2015-06-11T14:51:36-05:00","created":"2015-03-13T14:51:36-05:00"}}`)
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	e, _ := time.Parse(time.RFC3339, "2015-06-11T14:51:36-05:00")

	tkn, err := client.NewToken("test", "", "", e)
	if err != nil {
		t.Fatal(err)
	}

	if tkn.Name != "test" {
		t.Errorf("Name doesn't match, expected test, got %s", tkn.Name)
	}
	if c, _ := time.Parse(time.RFC3339, "2015-03-13T14:51:36-05:00"); !tkn.CreatedTime().Equal(c) {
		t.Errorf("Created time doesn't match for test token. Expected %v got %v", c, tkn.CreatedTime())
	}
	if e, _ := time.Parse(time.RFC3339, "2015-06-11T14:51:36-05:00"); !tkn.ExpiresTime().Equal(e) {
		t.Errorf("Expired time doesn't match for test token. Expected %v got %v", e, tkn.ExpiresTime())
	}
	if tkn.ID != "61c78f8bf7f62f2b17aaaa4e1345e27b" {
		t.Errorf("ID doesn't match, expected 61c78f8bf7f62f2b17aaaa4e1345e27b, got %s", tkn.ID)
	}

}

func TestDeleteToken(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/token/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" {
				fmt.Fprint(w, `{"status":"success"}`)
			} else {
				fmt.Fprint(w, `{"status":"success","data":
				{"id":"61c78f8bf7f62f2b17aaaa4e1345e27b","name":"test","expires":"2015-06-11T14:51:36-05:00","created":"2015-03-13T14:51:36-05:00"}}`)

			}
		})

	client, err := New(server.URL, username, password, nil)
	if err != nil {
		t.Fatal(err)
	}

	tkn, err := client.GetToken("")
	if err != nil {
		t.Fatal(err)
	}

	err = tkn.Delete()
	if err != nil {
		t.Fatalf("Error deleting token: %v", err)
	}

}
