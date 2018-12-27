package zgok

import (
	"testing"
)

func TestNewLocalFileSystem(t *testing.T) {
	zfs, err := NewLocalFileSystem("./testdata")
	if err != nil {
		t.Fatalf("Failed to add local file: %v", err)
	}

	f, err := zfs.ReadFileString("testdata/foo")
	if err != nil {
		t.Fatalf("Failed to testdata/foo read: %v", err)
	}
	if f != "foo" {
		t.Errorf("Expected 'foo' from testdata/foo but got: %v", f)
	}

	f, err = zfs.ReadFileString("testdata/dir/baz")
	if err != nil {
		t.Fatalf("Failed to testdata/dir/baz read: %v", err)
	}
	if f != "baz" {
		t.Errorf("Expected 'baz' from testdata/dir/baz but got: %v", f)
	}
}
