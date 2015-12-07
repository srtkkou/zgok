package zgok

import (
	"testing"
)

func TestZipUnzip(t *testing.T) {
	zipper := NewZipper()
	// Add file.
	err := zipper.Add("testdata/foo")
	if err != nil {
		t.Errorf("Failed to add file: %v", err)
	}
	// Add directory.
	err = zipper.Add("testdata/dir")
	if err != nil {
		t.Errorf("Failed to add directory: %v", err)
	}
	// Add empty directory.
	err = zipper.Add("testdata/empty")
	zipper.Close()
	bytes, err := zipper.Bytes()
	if err != nil {
		t.Errorf("Failed to get bytes: %v", err)
	}
	// Unzip files.
	unzipper := NewUnzipper(&bytes)
	zfs, err := unzipper.Unzip()
	if err != nil {
		t.Errorf("Failed to unzip files: %v", err)
	}
	// Verify paths.
	zfsPaths := zfs.Paths()
	expectedPaths := []string{
		"testdata/dir/bar",
		"testdata/dir/baz",
		"testdata/foo",
	}
	for i, expected := range expectedPaths {
		zfsPath := zfsPaths[i]
		if zfsPath != expected {
			t.Errorf("Paths(): expected [%s] got [%s]", expected, zfsPath)
		}
	}
	// Verify foo file.
	fooStr, err := zfs.ReadFileString("testdata/foo")
	if err != nil {
		t.Errorf("Failed to get [testdata/foo]: %v", err)
	}
	if fooStr != "foo" {
		t.Errorf("[testdata/foo]: expected [%s] got [%s]", "foo", fooStr)
	}
	// Verify bar file.
	barStr, err := zfs.ReadFileString("testdata/dir/bar")
	if err != nil {
		t.Errorf("Failed to get [testdata/dir/bar]: %v", err)
	}
	if barStr != "bar" {
		t.Errorf("[testdata/dir/bar]: expected [%s] got [%s]", "bar", barStr)
	}
	// Verify baz file.
	bazStr, err := zfs.ReadFileString("testdata/dir/baz")
	if err != nil {
		t.Errorf("Failed to get [testdata/dir/baz]: %v", err)
	}
	if bazStr != "baz" {
		t.Errorf("[testdata/dir/baz]: expected [%s] got [%s]", "baz", bazStr)
	}
	// Verify ignoring empty directory.
	_, err = zfs.GetFile("testdata/empty")
	if err == nil {
		t.Errorf("Empty directory [testdata/empty] is not ignored.")
	}
}
