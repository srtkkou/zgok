package zgok

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	APP_BYTE_SIZE       = 8
	SIGNATURE_BYTE_SIZE = 64
)

var (
	// byte order of the signature
	byteOrder binary.ByteOrder = binary.BigEndian
)

// Signature interface.
type Signature interface {
	ExeSize() int64
	SetExeSize(exeSize int64)
	ZipSize() int64
	SetZipSize(zipSize int64)
	TotalSize() int64
	String() string
	Dump() ([]byte, error)
}

// signature
type zgokSignature struct {
	app     string
	major   uint16
	minor   uint16
	rev     uint16
	exeSize int64
	zipSize int64
}

// Initialize signature.
func NewSignature() Signature {
	return &zgokSignature{
		app:   APP,
		major: MAJOR,
		minor: MINOR,
		rev:   REV,
	}
}

// Restore signature from bytes.
func RestoreSignature(data []byte) (Signature, error) {
	// Check size.
	if len(data) != SIGNATURE_BYTE_SIZE {
		return nil, errors.New("Invalid signature size.")
	}
	// Convert bytes to buffer.
	buf := bytes.NewBuffer(data)
	// Initialize signature.
	s := &zgokSignature{}
	// Restore app name.
	var appBytes []byte = make([]byte, APP_BYTE_SIZE, APP_BYTE_SIZE)
	_, err := buf.Read(appBytes)
	if err != nil {
		return nil, err
	}
	app, err := restoreAppString(appBytes)
	if err != nil {
		return nil, err
	}
	if app != APP {
		return nil, errors.New("Invalid signature.")
	}
	s.app = app
	// Restore major version.
	var major uint16
	err = binary.Read(buf, byteOrder, &major)
	if err != nil {
		return nil, err
	}
	s.major = major
	// Restore minor version.
	var minor uint16
	err = binary.Read(buf, byteOrder, &minor)
	if err != nil {
		return nil, err
	}
	s.minor = minor
	// Restore revision.
	var rev uint16
	err = binary.Read(buf, byteOrder, &rev)
	if err != nil {
		return nil, err
	}
	s.rev = rev
	// Restore exe size.
	var exeSize int64
	err = binary.Read(buf, byteOrder, &exeSize)
	if err != nil {
		return nil, err
	}
	if exeSize <= 0 {
		return nil, errors.New("Invalid signature.")
	}
	s.exeSize = exeSize
	// Restore zip size.
	var zipSize int64
	err = binary.Read(buf, byteOrder, &zipSize)
	if err != nil {
		return nil, err
	}
	if zipSize <= 0 {
		return nil, errors.New("Invalid signature.")
	}
	s.zipSize = zipSize
	return s, nil
}

// Restore app string from bytes.
func restoreAppString(appBytes []byte) (string, error) {
	// Check byte size.
	if len(appBytes) != APP_BYTE_SIZE {
		return "", errors.New("Invalid app byte size.")
	}
	// Get string length.
	appLen := bytes.IndexByte(appBytes, 0)
	if appLen < 0 || APP_BYTE_SIZE < appLen {
		appLen = APP_BYTE_SIZE
	}
	app := string(appBytes[:appLen])
	return app, nil
}

// Get exe file byte size.
func (sig *zgokSignature) ExeSize() int64 {
	return sig.exeSize
}

// Set exe file byte size.
func (sig *zgokSignature) SetExeSize(exeSize int64) {
	sig.exeSize = exeSize
}

// Get zip file byte size.
func (sig *zgokSignature) ZipSize() int64 {
	return sig.zipSize
}

// Set zip file byte size.
func (sig *zgokSignature) SetZipSize(zipSize int64) {
	sig.zipSize = zipSize
}

// Calculate the total byte size.
func (s *zgokSignature) TotalSize() int64 {
	return s.exeSize + s.zipSize + SIGNATURE_BYTE_SIZE
}

// Convert to string.
func (s *zgokSignature) String() string {
	return fmt.Sprintf("%s(exe:%d,zip:%d,total:%d)",
		Version(), s.exeSize, s.zipSize, s.TotalSize())
}

// Dump signature to bytes.
func (s *zgokSignature) Dump() ([]byte, error) {
	// Initialize buffer and byte count.
	buf := new(bytes.Buffer)
	byteCount := 0
	// Write app name.
	appBytes := s.appBytes()
	err := binary.Write(buf, byteOrder, appBytes)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(appBytes)
	// Write versions.
	err = binary.Write(buf, byteOrder, s.major)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.major)
	err = binary.Write(buf, byteOrder, s.minor)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.minor)
	err = binary.Write(buf, byteOrder, s.rev)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.rev)
	// Write sizes.
	err = binary.Write(buf, byteOrder, s.exeSize)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.exeSize)
	err = binary.Write(buf, byteOrder, s.zipSize)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.zipSize)
	// Fill with blank bytes.
	for i := byteCount; i < SIGNATURE_BYTE_SIZE; i++ {
		err := binary.Write(buf, byteOrder, byte(0))
		if err != nil {
			return []byte{}, err
		}
	}
	return buf.Bytes(), nil
}

// Byte array of the app string.
func (s *zgokSignature) appBytes() [APP_BYTE_SIZE]byte {
	var result [APP_BYTE_SIZE]byte
	appBytes := []byte(s.app)
	for i, _ := range result {
		if i < len(appBytes) {
			result[i] = appBytes[i]
		}
	}
	return result
}
