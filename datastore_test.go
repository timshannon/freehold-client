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

func TestDSGetValue(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/properties/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":{"name":"610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"url":"/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"permissions":{"owner":"tshannon","private":"rw"},"size":1048576,"modified":"2015-04-15T10:21:46-05:00"}}`)
		})

	mux.HandleFunc("/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				fmt.Fprint(w, `{"status":"success","data":{"name":"610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"url":"/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds"}}`)
			}

			if r.Method == "PUT" {
				fmt.Fprint(w, `{"status":"success"}`)
			}
			if r.Method == "GET" {
				fmt.Fprint(w, `{"status":"success","data":"testvalue"}`)
				return
			}
			if r.Method == "DELETE" {
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	client, err := New(server.URL, username, password)
	//client, err := New("https://tshannon.org", "tshannon", "_hvlkuhuCsGifYbxVBuEgaLbYvAb9kVz6dHA43XThCk=")
	if err != nil {
		t.Fatal(err)
	}

	ds, err := client.NewDatastore("/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds")
	if err != nil {
		t.Fatal(err)
	}
	err = ds.Put(10, "testvalue")
	if err != nil {
		t.Fatal(err)
	}

	result := ""
	err = ds.Get(10, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "testvalue" {
		t.Errorf("Expected testvalue, got %s", result)
	}

	err = ds.Delete(10)
	if err != nil {
		t.Fatal(err)
	}

	err = ds.Drop()
	if err != nil {
		t.Fatal(err)
	}

}

func TestDSIter(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/properties/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":{"name":"610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"url":"/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"permissions":{"owner":"tshannon","private":"rw"},"size":1048576,"modified":"2015-04-15T10:21:46-05:00"}}`)
		})

	mux.HandleFunc("/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				fmt.Fprint(w, `{"status":"success","data":{"name":"610c247c-520b-2d57-65d7-441b67a0b69d.ds",
					"url":"/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds"}}`)
			}

			if r.Method == "PUT" {
				fmt.Fprint(w, `{"status":"success"}`)
			}
			if r.Method == "GET" {
				body := requestBody(t, r)
				if strings.Contains(body, "min") {
					fmt.Fprint(w, `{"status":"success","data":{"key":10,"value":"minvalue"}}`)
					return
				}
				if strings.Contains(body, "max") {
					fmt.Fprint(w, `{"status":"success","data":{"key":100,"value":"maxvalue"}}`)
					return
				}
				if strings.Contains(body, "iter") {
					fmt.Fprint(w, `{"status":"success","data":[{"key":10,"value":"minvalue"}, {"key":100,"value":"maxvalue"}]}`)
					return
				}
			}
			if r.Method == "DELETE" {
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	client, err := New(server.URL, username, password)
	//client, err := New("https://tshannon.org", "tshannon", "_hvlkuhuCsGifYbxVBuEgaLbYvAb9kVz6dHA43XThCk=")
	if err != nil {
		t.Fatal(err)
	}

	ds, err := client.NewDatastore("/v1/datastore/testdata/610c247c-520b-2d57-65d7-441b67a0b69d.ds")
	if err != nil {
		t.Fatal(err)
	}
	err = ds.Put(10, "minvalue")
	if err != nil {
		t.Fatal(err)
	}

	err = ds.Put(100, "maxvalue")
	if err != nil {
		t.Fatal(err)
	}
	result := ""

	err = ds.Min().Value(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "minvalue" {
		t.Errorf("Expected minvalue, got %s", result)
	}

	err = ds.Max().Value(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "maxvalue" {
		t.Errorf("Expected maxvalue, got %s", result)
	}

	kvSlice, err := ds.Iter(&Iter{})
	if err != nil {
		t.Fatal(err)
	}

	err = kvSlice[0].Value(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "minvalue" {
		t.Errorf("Expected minvalue, got %s", result)
	}

	err = kvSlice[1].Value(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "maxvalue" {
		t.Errorf("Expected maxvalue, got %s", result)
	}

	err = ds.Drop()
	if err != nil {
		t.Fatal(err)
	}

}
