package zgok

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	APP   = "zgok"
	MAJOR = 0
	MINOR = 0
	REV   = 1
)

type zgok struct {
	isZgokFormat bool
	signature    *signature
	fileMap      map[string][]byte
}

func Initialize() *zgok {
	z := &zgok{}
	z.isZgokFormat = false
	// Get signature from the executable file itself.
	selfPath, _ := filepath.Abs(os.Args[0])
	signature, _ := parseSignature(selfPath)
	z.signature = signature
	// Setup file map.
	z.fileMap = make(map[string][]byte)
	if signature != nil {
		parseZip(signature, selfPath)
	}
	return z
}

func parseSignature(path string) (*signature, error) {
	// Get the size of the file.
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	totalSize := fileInfo.Size()
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Read signature bytes.
	buf := make([]byte, SIGNATURE_BYTE_SIZE)
	_, err = file.ReadAt(buf, totalSize-SIGNATURE_BYTE_SIZE)
	if err != nil {
		return nil, err
	}
	// Parse signature.
	signature, err := RestoreSignature(buf)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func parseZip(signature *signature, path string) error {
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Get zip section.
	exeSize := signature.exeSize
	zipSize := signature.zipSize
	zipSectionReader := io.NewSectionReader(file, exeSize, zipSize)
	zipReader, _ := zip.NewReader(zipSectionReader, zipSize)
	var rc io.ReadCloser
	for _, f := range zipReader.File {
		// Get file info.
		fileInfo := f.FileHeader.FileInfo()
		fmt.Printf("path: %s size: %d dir: %v\n", f.Name, fileInfo.Size(), fileInfo.IsDir())
		// Store bytes.
		buf := new(bytes.Buffer)
		rc, err = f.Open()
		if err != nil {
			fmt.Errorf("%v", err)
		}
		_, err = io.Copy(buf, rc)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		fmt.Println(len(buf.Bytes()))
		rc.Close()
	}
	return nil
}

func (z *zgok) ReadFile(path string) ([]byte, error) {
	return []byte{}, nil
}
