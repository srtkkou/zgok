package zgok

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Files interface {
	ReadFile(path string) ([]byte, error)
}

type zgokFiles struct {
	isZgokFormat bool
	signature    *signature
	fileMap      map[string][]byte
}

func Initialize() Files {
	z := &zgokFiles{}
	z.isZgokFormat = false
	// Get signature from the executable file itself.
	selfPath, _ := filepath.Abs(os.Args[0])
	fmt.Println("before signature init")
	signature, err := parseSignature(selfPath)
	if err != nil {
		fmt.Println(err)
	}
	z.signature = signature
	if z.signature != nil {
		fmt.Println(z.signature.String())
	}

	// Setup file map.
	z.fileMap = make(map[string][]byte)
	z.parseZip(selfPath)

	// Show file map.
	fmt.Println("--FileMap--")
	for path, _ := range z.fileMap {
		fmt.Printf("key=%s\n", path)
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

func (z *zgokFiles) parseZip(path string) error {
	// Do nothing if signature is nil.
	if z.signature == nil {
		return nil
	}
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Get zip section.
	exeSize := z.signature.exeSize
	zipSize := z.signature.zipSize
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
			fmt.Printf("%v", err)
		}
		_, err = io.Copy(buf, rc)
		if err != nil {
			fmt.Printf("%v", err)
		}
		// Store into file map.
		z.fileMap[f.Name] = buf.Bytes()
		fmt.Println(len(buf.Bytes()))
		rc.Close()
	}
	return nil
}

func (z *zgokFiles) ReadFile(path string) ([]byte, error) {
	// Read from file system if signature is blank.
	fmt.Printf("path:%s\n", path)
	if z.signature == nil {
		fmt.Println("READ FROM FILE SYSTEM.")
		return ioutil.ReadFile(path)
	}
	// Read file from zip.
	zipPath := filepath.Join(APP, path)
	fmt.Printf("READ FROM ZIP.%s\n", zipPath)
	bytes, exists := z.fileMap[zipPath]
	if !exists {
		return []byte{}, errors.New("File does not exist.")
	}
	fmt.Printf("content=%s\n", string(bytes))
	return bytes, nil
}
