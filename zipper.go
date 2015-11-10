package zgok

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Zipper struct {
	isClosed bool
	buffer   *bytes.Buffer
	writer   *zip.Writer
	basePath string
}

func NewZipper() *Zipper {
	z := &Zipper{}
	z.isClosed = false
	z.basePath = "zgok"
	z.buffer = new(bytes.Buffer)
	z.writer = zip.NewWriter(z.buffer)
	return z
}

func (z *Zipper) Add(path string) error {
	// Check if zip is closed or not.
	if z.isClosed {
		return errors.New("ZIP is already closed.")
	}
	// Get file information.
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Determine if the path is a directory or not.
	if fileInfo.IsDir() {
		err = z.addDir(path)
	} else {
		err = z.addFile(path)
	}
	if err != nil {
		return err
	}
	return nil
}

func (z *Zipper) Close() error {
	err := z.writer.Close()
	if err != nil {
		return err
	}
	z.isClosed = true
	return nil
}

func (z *Zipper) Bytes() ([]byte, error) {
	if !z.isClosed {
		return []byte{}, errors.New("ZIP is not closed.")
	}
	return z.buffer.Bytes(), nil
}

func (z *Zipper) addFile(filePath string) error {
	// Get file information.
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	// Get file content.
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Set zip header.
	header, _ := zip.FileInfoHeader(fileInfo)
	path := filepath.Join(z.basePath, filePath)
	header.Name = path
	zipFile, err := z.writer.CreateHeader(header)
	if err != nil {
		return err
	}
	// Write content.
	_, err = zipFile.Write(content)
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

func (z *Zipper) addDir(dirPath string) error {
	// List files in directory.
	pattern := filepath.Join(dirPath, "**", "*")
	filePaths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	// Add all files to zip.
	for _, filePath := range filePaths {
		err = z.addFile(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
