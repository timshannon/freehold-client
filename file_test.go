package freeholdclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	username = "tester"
	password = "testerToken"
	dirPath  = "/v1/file/testing/"
	filePath = "/v1/file/testing/test.txt"
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

	f, err := client.GetFile(dirPath)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(f.Permissions.Owner)

}

// will go away and be replaced with legitimate standalone tests with a mock server
func TestRealFile(t *testing.T) {
	client, err := New("https://dev.tshannon.org", "tshannon", "zVZm8Ic0_zGUpxTbry9Ph4C--0vs4v9-QM6uYITXk4g=",
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		t.Fatal(err)
	}

	f, err := client.GetFile(dirPath)

	if err != nil {
		t.Fatal(err)
	}
	//children, err := f.Children()
	//if err != nil {
	//t.Fatal(err)
	//}
	//for i := range children {
	//t.Log(children[i].URL)
	//}

	localFile, err := os.Open("test.txt")
	defer localFile.Close()

	if err != nil {
		t.Fatal(err)
	}

	upFile, err := client.UploadFile(localFile, f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Upload Successful: ", upFile.Name)

	//f2, err := client.GetFile(filePath)
	//if err != nil {
	//t.Fatal(err)
	//}

	//lf, err := os.Create(f2.Name)
	//defer lf.Close()
	//if err != nil {
	//t.Fatal(err)

	//}

	//_, err = io.Copy(lf, f2)
	//defer f2.Close()
	//if err != nil {
	//t.Fatal(err)
	//}
}
