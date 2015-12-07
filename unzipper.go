package zgok

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
)

// Unzipper.
type Unzipper struct {
	isUnzipped bool          // Is the file already unzipped?
	reader     *bytes.Reader // Byte reader.
	size       int64         // Size of the zipped file.
}

// Create new unzipper.
func NewUnzipper(zipBytes *[]byte) *Unzipper {
	u := &Unzipper{
		isUnzipped: false,
		reader:     bytes.NewReader(*zipBytes),
		size:       int64(len(*zipBytes)),
	}
	return u
}

// Unzip all the files in zip.
func (u *Unzipper) Unzip() (FileSystem, error) {
	var err error
	// Check if it is already unzipped.
	if u.isUnzipped {
		return nil, errors.New("Already unzipped.")
	}
	// Initialize zip reader.
	zipReader, err := zip.NewReader(u.reader, u.size)
	if err != nil {
		return nil, err
	}
	// Prepare file system.
	zfs := NewFileSystem()
	// Get all files.
	var readCloser io.ReadCloser
	for _, file := range zipReader.File {
		// Initialize zgok file.
		zgokFile := NewZgokFile()
		// Set file path.
		path := file.FileHeader.Name
		zgokFile.SetPath(path)
		// Set file info.
		fileInfo := file.FileHeader.FileInfo()
		zgokFile.SetFileInfo(fileInfo)
		// Open file.
		readCloser, err = file.Open()
		if err != nil {
			break
		}
		// Copy bytes.
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, readCloser)
		if err != nil {
			break
		}
		zgokFile.SetBytes(buf.Bytes())
		// Close file.
		readCloser.Close()
		// Add file to file system.
		zfs.AddFile(zgokFile)
	}
	if err != nil {
		return nil, err
	}
	u.isUnzipped = true
	return zfs, nil
}
