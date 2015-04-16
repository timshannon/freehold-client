// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":{"type":"basic","user":"tester","homeApp":"home","admin":true}}`)
		})

	client, err := New(server.URL, username, password)
	//client, err := New("https://tshannon.org", "tshannon", "_hvlkuhuCsGifYbxVBuEgaLbYvAb9kVz6dHA43XThCk=")
	if err != nil {
		t.Fatal(err)
	}

	ath, err := client.Auth()
	if err != nil {
		t.Fatal(err)
	}

	if ath.Username != username {
		t.Errorf("Username doesn't match, expected tester, got %s", ath.Username)
	}
}
