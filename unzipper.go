package zgok

import (
	"archive/zip"
)

type Unzipper struct {
	bytes  []byte
	reader *zip.Reader
}

func NewUnzipper(bytes []byte) *Unzipper {
	u := &Unzipper{}
	return u
}
