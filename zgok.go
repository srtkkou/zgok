/*
Package zgok provides a simple library to create a single binary with asset files in Go (Golang).

Example:

	package main

	import (
		"net/http"
		"github.com/srtkkou/zgok"
		"os"
	)

	func main() {
		zfs, err := zgok.RestoreFileSystem(os.Args[0])
		if err != nil {
			panic(err)
		}
		assetServer := zfs.FileServer("web/public")
		http.Handle("/assets/", http.StripPrefix("/assets/", assetServer))
		http.ListenAndServe(":8080", nil)
	}
*/
package zgok

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	APP   = "zgok" // Application name.
	MAJOR = 0      // Major version.
	MINOR = 0      // Minor version.
	REV   = 1      // Revision.
)

// Get version string.
func Version() string {
	return fmt.Sprintf("%s-%d.%d.%d", APP, MAJOR, MINOR, REV)
}

// File system interface.
// Implements [net/http.FileSystem]
type FileSystem interface {
	AddFile(file File)
	GetFile(path string) (File, error)
	ReadFile(path string) ([]byte, error)
	ReadFileString(path string) (string, error)
	Paths() []string
	SubFileSystem(rootPath string) (FileSystem, error)
	Signature() Signature
	SetSignature(signature Signature)
	String() string
	Open(name string) (http.File, error)     // Implements [net/http.FileSystem.Open]
	FileServer(basePath string) http.Handler // Get a static file server.
}

// Zgok file system.
type zgokFileSystem struct {
	signature Signature       // Zgok signature.
	rootPath  string          // Root path of the file system.
	fileMap   map[string]File // Map of files.
}

// Create a new file system.
func NewFileSystem() FileSystem {
	return &zgokFileSystem{
		signature: nil,
		rootPath:  APP,
		fileMap:   make(map[string]File),
	}
}

// Restore file system.
func RestoreFileSystem(path string) (FileSystem, error) {
	// Get bytes of exe file.
	exeBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Restore signature.
	sigOffset := len(exeBytes) - SIGNATURE_BYTE_SIZE
	sigBytes := exeBytes[sigOffset:]
	signature, err := RestoreSignature(sigBytes)
	if err != nil {
		return nil, err
	}
	// Unzip zip section.
	zipOffset := signature.ExeSize()
	zipBytes := exeBytes[zipOffset:sigOffset]
	unzipper := NewUnzipper(&zipBytes)
	zfs, err := unzipper.Unzip()
	if err != nil {
		return nil, err
	}
	// Set signature.
	zfs.SetSignature(signature)
	return zfs, nil
}

// Add file to file system.
func (zfs *zgokFileSystem) AddFile(file File) {
	key := filepath.ToSlash(file.Path())
	zfs.fileMap[key] = file
}

// Get file from file system.
func (zfs *zgokFileSystem) GetFile(path string) (File, error) {
	key := filepath.ToSlash(filepath.Join(zfs.rootPath, path))
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
		return []byte{}, err
	}
	return file.Bytes(), nil
}

// Get the content of file in string from file system.
func (zfs *zgokFileSystem) ReadFileString(path string) (string, error) {
	bytes, err := zfs.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(bytes)
	return str, nil
}

// Get all the paths stored in the file system.
func (zfs *zgokFileSystem) Paths() []string {
	paths := []string{}
	prefix := zfs.rootPath + "/"
	for key := range zfs.fileMap {
		if strings.HasPrefix(key, prefix) {
			relPath := strings.TrimPrefix(key, prefix)
			paths = append(paths, relPath)
		}
	}
	sort.Strings(paths)
	return paths
}

// Get a sub file system.
func (zfs *zgokFileSystem) SubFileSystem(rootPath string) (FileSystem, error) {
	// Check root path.
	if strings.Contains(rootPath, "..") {
		return nil, errors.New("No double dots [..] are allowed in root path.")
	}
	// Initialize sub file system.
	newRootPath := filepath.ToSlash(filepath.Join(zfs.rootPath, rootPath))
	subFs := &zgokFileSystem{
		signature: zfs.signature,
		rootPath:  newRootPath,
		fileMap:   make(map[string]File),
	}
	// Add all the sets matching the new root path.
	for key, value := range zfs.fileMap {
		if key[0:len(newRootPath)] == newRootPath {
			subFs.fileMap[key] = value
		}
	}
	return subFs, nil
}

// Get signature.
func (zfs *zgokFileSystem) Signature() Signature {
	return zfs.signature
}

// Set signature.
func (zfs *zgokFileSystem) SetSignature(signature Signature) {
	zfs.signature = signature
}

// Get string.
func (zfs *zgokFileSystem) String() string {
	return zfs.Signature().String()
}

// Open the file.
// Implements [net/http.FileSystem.Open]
func (zfs *zgokFileSystem) Open(name string) (http.File, error) {
	path := strings.TrimLeft(name, "/")
	file, err := zfs.GetFile(path)
	if err != nil {
		// Return an abstract directory.
		dir := &zgokFile{
			fileInfo: zgokFileInfo{
				name: filepath.Base(path),
				mode: os.ModeDir | os.ModePerm,
			},
		}
		return dir, nil
	}
	// Set a new file reader
	file.SetNewReader()
	return file, nil
}

// File interface.
type File interface {
	SetPath(path string)                          // Set file path.
	Path() string                                 // Get file path.
	SetFileInfo(fileInfo os.FileInfo)             // Set file info.
	FileInfo() os.FileInfo                        // Get file info.
	SetBytes(content []byte)                      // Set content bytes.
	Bytes() []byte                                // Get content bytes.
	SetNewReader()                                // Set a new reader.
	Close() error                                 // Implements [net/http.File.Close]
	Read(p []byte) (int, error)                   // Implements [net/http.File.Read]
	Readdir(count int) ([]os.FileInfo, error)     // Implements [net/http.File.Readdir]
	Seek(offset int64, whence int) (int64, error) // Implements [net/http.File.Seek]
	Stat() (os.FileInfo, error)                   // Implements [net/http.File.Stat]
}

// Zgok file.
type zgokFile struct {
	path     string        // Path of the file.
	fileInfo os.FileInfo   // File info.
	content  []byte        // Content of the file.
	reader   *bytes.Reader // File reader.
}

// Create a new zgok file.
func NewZgokFile() File {
	return &zgokFile{}
}

// Set file path.
func (zf *zgokFile) SetPath(path string) {
	zf.path = strings.Replace(path, `\`, "/", -1)
}

// Get file path.
func (zf *zgokFile) Path() string {
	return zf.path
}

// Set file info.
func (zf *zgokFile) SetFileInfo(fileInfo os.FileInfo) {
	zf.fileInfo = fileInfo
}

// Get file info.
func (zf *zgokFile) FileInfo() os.FileInfo {
	return zf.fileInfo
}

// Set content bytes.
func (zf *zgokFile) SetBytes(content []byte) {
	zf.content = content
}

// Get content bytes.
func (zf *zgokFile) Bytes() []byte {
	return zf.content
}

// Set a new reader.
func (zf *zgokFile) SetNewReader() {
	reader := bytes.NewReader(zf.content)
	zf.reader = reader
}

// Close file.
// Implements [net/http.File.Close]
func (zf *zgokFile) Close() error {
	return nil
}

// Read file.
// Implements [net/http.File.Read]
func (zf *zgokFile) Read(p []byte) (int, error) {
	return zf.reader.Read(p)
}

// Read directories.
// Implements [net/http.File.Readdir]
func (zf *zgokFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("Readdir is not allowed.")
}

// Seek file.
// Implements [net/http.File.Seek]
func (zf *zgokFile) Seek(offset int64, whence int) (int64, error) {
	return zf.reader.Seek(offset, whence)
}

// Get file info.
// Implements [net/http.File.Stat]
func (zf *zgokFile) Stat() (os.FileInfo, error) {
	return zf.fileInfo, nil
}

// Zgok file info.
// Implements [os.FileInfo]
type zgokFileInfo struct {
	name    string      // File name.
	size    int64       // Size of the file.
	mode    os.FileMode // File mode.
	modTime time.Time   // Modified time of the file.
}

// Get name.
// Implements [os.FileInfo.Name]
func (i zgokFileInfo) Name() string {
	return i.name
}

// Get size.
// Implements [os.FileInfo.Size]
func (i zgokFileInfo) Size() int64 {
	return i.size
}

// Get mode.
// Implements [os.FileInfo.Mode]
func (i zgokFileInfo) Mode() os.FileMode {
	return i.mode
}

// Get modified time.
// Implements [os.FileInfo.ModTime]
func (i zgokFileInfo) ModTime() time.Time {
	return i.modTime
}

// Check if it is a directory.
// Implements [os.FileInfo.IsDir]
func (i zgokFileInfo) IsDir() bool {
	return i.mode.IsDir()
}

// Get sys information. (Only returns nil.)
// Implements [os.FileInfo.Sys]
func (i zgokFileInfo) Sys() interface{} {
	return nil
}
