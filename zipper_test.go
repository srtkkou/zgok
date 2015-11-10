package zgok

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestZipper(t *testing.T) {
	zipper := NewZipper()
	err := zipper.Add("zipper_test.go")
	if err != nil {
		t.Errorf("%v", err)
	}
	err = zipper.Add("cmd")
	if err != nil {
		t.Errorf("%v", err)
	}
	zipper.Close()
	data, _ := zipper.Bytes()
	err = ioutil.WriteFile("zip_path.zip", data, os.ModePerm)
	if err != nil {
		t.Errorf("%v", err)
	}
}
