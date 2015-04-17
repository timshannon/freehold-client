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

func TestBackups(t *testing.T) {
	startMockServer()
	defer stopMockServer()

	//Setup Mock Handler
	mux.HandleFunc("/v1/backup/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fmt.Fprint(w, `{"status":"success","data":[{"when":"2015-04-15T10:07:18-05:00","file":"/v1/file/backup_test-2ddc9e35-b1de-233a-a092-5f2435697650.zip","who":"quinitTestUserAdmin","datastores":["settings.ds","app.ds","log.ds","backup.ds","permission.ds","ratelimit.ds","token.ds","session.ds","user.ds"]},{"when":"2015-04-17T13:50:32-05:00","file":"/v1/file/backups/freehold_backup-2015-04-17T13:50:32-05:00.zip","who":"tshannon","datastores":["settings.ds","app.ds","log.ds","backup.ds","permission.ds","ratelimit.ds","token.ds","session.ds","user.ds"]}]}`)
			}

			if r.Method == "POST" {
				fmt.Fprint(w, `{"status":"success","data":"/v1/file/backups/freehold_backup-2015-04-17T13:50:32-05:00.zip"}`)
			}
		})

	client, err := New(server.URL, username, password)
	if err != nil {
		t.Fatal(err)
	}

	file, err := client.NewBackup("", []string{})
	if err != nil {
		t.Fatal(err)
	}

	backups, err := client.GetBackups(time.Now().AddDate(0, 0, -4), time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for i := range backups {
		if backups[i].File == file {
			found = true
		}
	}

	if !found {
		t.Fatalf("Backup files %s not found.", file)
	}

}
