package zgok

import (
	"archive/zip"
	"bytes"
	"errors"
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

// Create new zipper.
func NewZipper() *Zipper {
	z := &Zipper{}
	z.isClosed = false
	z.basePath = "zgok"
	z.buffer = new(bytes.Buffer)
	z.writer = zip.NewWriter(z.buffer)
	return z
}

// Add files in the path to zip.
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

// Close zip writer.
func (z *Zipper) Close() error {
	err := z.writer.Close()
	if err != nil {
		return err
	}
	z.isClosed = true
	return nil
}

// Get bytes of zip.
func (z *Zipper) Bytes() ([]byte, error) {
	if !z.isClosed {
		return []byte{}, errors.New("ZIP is not closed.")
	}
	return z.buffer.Bytes(), nil
}

// Add file to zip.
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
	return nil
}

// Add directory to zip.
func (z *Zipper) addDir(dirPath string) error {
	// Walk through all the files in the directory.
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			// Do nothing on directory.
			if info.IsDir() {
				return nil
			}
			// Add file to zip.
			err = z.addFile(path)
			if err != nil {
				return err
			}
			return nil
		})
	return err
}
