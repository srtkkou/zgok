package zgok

import (
	"net/http"
)

// Get a static file server.
func (zfs *zgokFileSystem) FileServer(basePath string) http.Handler {
	var server http.Handler
	subFs, err := zfs.SubFileSystem(basePath)
	if err == nil {
		server = http.FileServer(subFs)
	}
	return server
}
