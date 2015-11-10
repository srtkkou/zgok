package zgok

import (
	"io/ioutil"
	"os"
	//"bytes"
	"testing"
	//"archive/zip"
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
	//t.Errorf("%v", data[0:10])
	err = ioutil.WriteFile("zip_path.zip", data, os.ModePerm)
	if err != nil {
		t.Errorf("%v", err)
	}
}
