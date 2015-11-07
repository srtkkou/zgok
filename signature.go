package zgok

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	ID                  = "zgok"
	ID_BYTE_SIZE        = 8
	MAJOR_VERSION       = 0
	MINOR_VERSION       = 1
	SIGNATURE_BYTE_SIZE = 64
)

var (
	// byte order of the signature
	byteOrder binary.ByteOrder = binary.BigEndian
)

type SignatureInterface interface {
	//Verify() error

	Bytes() []byte
	Id() string
	MajorVersion() uint16
	MinorVersion() uint16
	ExeSize() int64
	ZipSize() int64
	String() string
	SetExeSize(exeSize int64)
	SetZipSize(zipSize int64)
}

type signature struct {
	id           string
	majorVersion uint16
	minorVersion uint16
	exeSize      int64
	zipSize      int64
}

// Initialize signature.
func NewSignature() *signature {
	return &signature{
		id:           ID,
		majorVersion: MAJOR_VERSION,
		minorVersion: MINOR_VERSION,
	}
}

func RestoreSignature(bytes []byte) (*signature, error) {
	// Check size.
	if len(bytes) != SIGNATURE_BYTE_SIZE {
		return nil, errors.New("Invalid signature size.")
	}
	// Initialize signature.
	s := &signature{}
	// Restore id.
	//idBytes :=
	return s, nil
}

func (sig *signature) SetExeSize(exeSize int64) {
	sig.exeSize = exeSize
}

func (sig *signature) SetZipSize(zipSize int64) {
	sig.zipSize = zipSize
}

func (s *signature) TotalSize() int64 {
	return s.exeSize + s.zipSize + SIGNATURE_BYTE_SIZE
}

func (s *signature) String() string {
	return fmt.Sprintf("%s-%d.%d(exe:%d,zip:%d,total:%d)",
		s.id, s.majorVersion, s.minorVersion,
		s.exeSize, s.zipSize, s.TotalSize())
}

// Dump signature to bytes.
func (s *signature) Dump() ([]byte, error) {
	// Initialize buffer and byte count.
	buf := new(bytes.Buffer)
	byteCount := 0
	// Copy ID.
	err := binary.Write(buf, byteOrder, s.idBytes())
	if err != nil {
		return []byte{}, errors.New("Failed to write id.")
	}
	byteCount += binary.Size(s.idBytes())
	// Copy version informations.
	err = binary.Write(buf, byteOrder, s.majorVersion)
	if err != nil {
		return []byte{}, errors.New("Failed to write major version.")
	}
	byteCount += binary.Size(s.majorVersion)
	err = binary.Write(buf, byteOrder, s.minorVersion)
	if err != nil {
		return []byte{}, errors.New("Failed to write major version.")
	}
	byteCount += binary.Size(s.minorVersion)
	// Copy size informations.
	err = binary.Write(buf, byteOrder, s.exeSize)
	if err != nil {
		return []byte{}, errors.New("Failed to write exe size.")
	}
	byteCount += binary.Size(s.exeSize)
	err = binary.Write(buf, byteOrder, s.zipSize)
	if err != nil {
		return []byte{}, errors.New("Failed to write zip size.")
	}
	byteCount += binary.Size(s.zipSize)
	// Fill with blank bytes.
	for i := byteCount; i < SIGNATURE_BYTE_SIZE; i++ {
		err := binary.Write(buf, byteOrder, byte(0))
		if err != nil {
			return []byte{}, errors.New("Failed to write blank bytes.")
		}
	}
	return buf.Bytes(), nil
}

// Byte array of the id string.
func (s *signature) idBytes() [ID_BYTE_SIZE]byte {
	var result [ID_BYTE_SIZE]byte
	idBytes := []byte(s.id)
	for i, _ := range result {
		if i < len(idBytes) {
			result[i] = idBytes[i]
		}
	}
	return result
}
