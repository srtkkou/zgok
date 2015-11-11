package zgok

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	APP   = "zgok"
	MAJOR = 0
	MINOR = 0
	REV   = 1
)

// File system interface.
type FileSystem interface {
	AddFile(file File)
	GetFile(path string) (File, error)
	ReadFile(path string) ([]byte, error)
	ReadFileString(path string) (string, error)
}

// Zgok file system.
type zgokFileSystem struct {
	fileMap map[string]File
}

// Create a new file system.
func NewFileSystem() FileSystem {
	return &zgokFileSystem{
		fileMap: make(map[string]File),
	}
}

// Restore file system.
func RestoreFileSystem() FileSystem {
	return nil
}

// Add file to file system.
func (zfs *zgokFileSystem) AddFile(file File) {
	key := zfs.toKey(file.FileInfo().Name())
	zfs.fileMap[key] = file
}

// Get file from file system.
func (zfs *zgokFileSystem) GetFile(path string) (File, error) {
	key := zfs.toKey(path)
	file, exists := zfs.fileMap[key]
	if !exists {
		return nil, errors.New("File doesn't exist.")
	}
	return file, nil
}

// Get the content of file in bytes from file system.
func (zfs *zgokFileSystem) ReadFile(path string) ([]byte, error) {
	file, err := zfs.GetFile(path)
	if err != nil {
		return []byte{}, nil
	}
	return file.Bytes(), nil
}

// Get the content of file in string from file system.
func (zfs *zgokFileSystem) ReadFileString(path string) (string, error) {
	bytes, err := zfs.ReadFile(path)
	if err != nil {
		return "", nil
	}
	str := string(bytes)
	return str, nil
}

// Convert path to zfs key.
func (zfs *zgokFileSystem) toKey(path string) string {
	// Convert to path with slash.
	key := filepath.ToSlash(path)
	// Strip app name from key if exists.
	if key[0:len(APP)] == APP {
		stripSize := len(APP) + 1
		key = key[stripSize:]
	}
	return key
}

// File interface.
type File interface {
	SetFileInfo(fileInfo os.FileInfo)
	FileInfo() os.FileInfo
	SetBytes(bytes []byte)
	Bytes() []byte
}

// Zgok file.
type zgokFile struct {
	fileInfo os.FileInfo
	content  []byte
}

// Create a new zgok file.
func NewZgokFile() File {
	return &zgokFile{}
}

// Set file info to file.
func (zf *zgokFile) SetFileInfo(fileInfo os.FileInfo) {
	zf.fileInfo = fileInfo
}

// Get file info of file.
func (zf *zgokFile) FileInfo() os.FileInfo {
	return zf.fileInfo
}

// Set bytes to file.
func (zf *zgokFile) SetBytes(bytes []byte) {
	zf.content = bytes
}

// Get bytes from file.
func (zf *zgokFile) Bytes() []byte {
	return zf.content
}
