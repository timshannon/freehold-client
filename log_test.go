// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLogs(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/log/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status":"success","data":[{"when":"2015-04-17T12:45:01-05:00","type":"server error","log":"http: TLS handshake error from 61.240.144.65:60000: read tcp 61.240.144.65:60000: connection reset by peer\n"},{"when":"2015-04-17T12:33:00-05:00","type":"server error","log":"http: TLS handshake error from 54.144.23.48:56029: remote error: unknown certificate authority\n"}]}`)
		})

	client, err := New(server.URL, username, password)
	if err != nil {
		t.Fatal(err)
	}

	logs, err := client.GetLogs(&LogIter{
		Type: "error",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) < 1 {
		t.Errorf("No Logs returned!")
	}
}
