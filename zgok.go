package zgok

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	APP   = "zgok"
	MAJOR = 0
	MINOR = 0
	REV   = 1
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
	SubFileSystem(rootPath string) (FileSystem, error)
	Signature() Signature
	SetSignature(signature Signature)
	String() string
	// Implements [net/http.FileSystem.Open]
	Open(name string) (http.File, error)
}

// Zgok file system.
type zgokFileSystem struct {
	signature Signature
	rootPath  string
	fileMap   map[string]File
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
func RestoreFileSystem() (FileSystem, error) {
	// Get bytes of exe file.
	exeBytes, err := ioutil.ReadFile(os.Args[0])
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
	key := filepath.ToSlash(file.FileInfo().Name())
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
	SetFileInfo(fileInfo os.FileInfo)
	FileInfo() os.FileInfo
	SetBytes(content []byte)
	Bytes() []byte
	SetNewReader()
	// Implements [net/http.File.Close]
	Close() error
	// Implements [net/http.File.Read]
	Read(p []byte) (int, error)
	// Implements [net/http.File.Readdir]
	Readdir(count int) ([]os.FileInfo, error)
	// Implements [net/http.File.Seek]
	Seek(offset int64, whence int) (int64, error)
	// Implements [net/http.File.Stat]
	Stat() (os.FileInfo, error)
}

// Zgok file.
type zgokFile struct {
	fileInfo os.FileInfo
	content  []byte
	reader   *bytes.Reader
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
func (zf *zgokFile) SetBytes(content []byte) {
	zf.content = content
}

// Get bytes from file.
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
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
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
