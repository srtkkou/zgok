package zgok

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// NewLocalFileSystem returns a FileSystem abstracted
// over the actual file system. This allows the RestoreFileSystem()
// to be used transparently with an automatic fallback
// to the local system when running during development.
func NewLocalFileSystem(rootPaths ...string) (FileSystem, error) {
	zfs := &zgokFileSystem{
		signature: nil,
		rootPath:  "",
		fileMap:   make(map[string]File),
	}

	for _, rp := range rootPaths {
		err := filepath.Walk(rp, func(p string, fInfo os.FileInfo, err error) error {
			zgokFile := NewZgokFile()
			zgokFile.SetPath(p)
			zgokFile.SetFileInfo(fInfo)

			// Read the file into memory.
			if !fInfo.IsDir() {
				f, err := os.Open(p)
				if err != nil {
					return err
				}

				// Copy bytes.
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, f)
				if err != nil {
					return err
				}
				zgokFile.SetBytes(buf.Bytes())

				// Add the file to the filesystem.
				zfs.AddFile(zgokFile)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return zfs, nil
}
