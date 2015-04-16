// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSession(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/auth/session/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fmt.Fprint(w, `{"status":"success","data":[{"id":"456FcIV6H7OxVeCF4SwHy+/z8Vhgu4qinq1HqC6qAxo=","CSRFToken":"jYTZjLHp4HRjDnwgKkSlkQcmOJxtSFQteUDbJblym08=","ipAddress":"11.111.111.111","created":"2015-04-16T11:50:47-05:00","userAgent":"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"},{"id":"J5q8BEyolZm6FOU99OOWR/htZq7w2UnlDSegHDxWJ94=","expires":"2015-05-01T01:58:28.020Z","CSRFToken":"lWEr1H94edR4ymhPoTVTtQHHBLNQ8hi2DAbtXFi9kSk=","ipAddress":"11.111.111.111","created":"2015-04-15T20:58:27-05:00","userAgent":"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"}]}`)
			}

			if r.Method == "DELETE" {
				// delete session
				fmt.Fprint(w, `{"status":"success"}`)
			}
		})

	client, err := New(server.URL, username, password)
	//client, err := New("https://tshannon.org", "tshannon", "_hvlkuhuCsGifYbxVBuEgaLbYvAb9kVz6dHA43XThCk=")
	if err != nil {
		t.Fatal(err)
	}

	s, err := client.AllSessions()
	if err != nil {
		t.Fatal(err)
	}

	if len(s) != 2 {
		t.Errorf("Session count doesn't match, expected 2, got %d", len(s))
	}
	for i := range s {
		err = s[i].Delete()
		if err != nil {
			t.Fatal(err)
		}
	}
}
