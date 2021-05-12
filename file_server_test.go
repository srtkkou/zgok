package zgok

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileServer(t *testing.T) {
	// Build zgok file.
	exePath := "testdata/executable"
	outPath := "file_server_test.out"
	builder := NewZgokBuilder()
	builder.SetExePath(exePath)
	builder.AddZipPath("testdata/foo")
	builder.AddZipPath("testdata/dir")
	builder.SetOutPath(outPath)
	builder.Build()
	// Load zgok filesystem.
	zfs, err := RestoreFileSystem(outPath)
	if err != nil {
		t.Errorf("RestoreFileSystem():error=[%v]", err)
	}
	// Initialize server.
	fileServer := zfs.FileServer("testdata")
	if fileServer == nil {
		t.Errorf("Failed to initialize file server.")
	}
	ts := httptest.NewServer(fileServer)
	defer ts.Close()
	// Get "foo"
	res, err := http.Get(ts.URL + "/foo")
	if err != nil {
		t.Errorf("Get [/foo] failed.")
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll():error=[%v]", err)
	}
	fooBody := string(content)
	if fooBody != "foo" {
		t.Errorf(`expected "foo" got "%v"`, fooBody)
	}
	// Get "dir/bar"
	res, err = http.Get(ts.URL + "/dir/bar")
	if err != nil {
		t.Errorf(`http.Get() failed to get "/dir/bar"`)
	}
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(`ioutil.ReadAll() %s`, err.Error())
	}
	barBody := string(content)
	if barBody != "bar" {
		t.Errorf(`expected "bar" got "%v"`, barBody)
	}
}
