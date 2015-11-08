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

// signature
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

// Restore signature from bytes.
func RestoreSignature(data []byte) (*signature, error) {
	// Check size.
	if len(data) != SIGNATURE_BYTE_SIZE {
		return nil, errors.New("Invalid signature size.")
	}
	// Convert bytes to buffer.
	buf := bytes.NewBuffer(data)
	// Initialize signature.
	s := &signature{}
	// Restore ID.
	var idBytes []byte = make([]byte, ID_BYTE_SIZE, ID_BYTE_SIZE)
	n, err := buf.Read(idBytes)
	if n != ID_BYTE_SIZE || err != nil {
		return nil, err
	}
	idLen := bytes.IndexByte(idBytes, 0)
	if idLen < 0 {
		idLen = ID_BYTE_SIZE
	}
	s.id = string(idBytes[:idLen])
	// Restore versions.
	err = binary.Read(buf, byteOrder, &s.majorVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, byteOrder, &s.minorVersion)
	if err != nil {
		return nil, err
	}
	// Restore sizes.
	err = binary.Read(buf, byteOrder, &s.exeSize)
	if s.exeSize <= 0 || err != nil {
		return nil, err
	}
	err = binary.Read(buf, byteOrder, &s.zipSize)
	if s.zipSize <= 0 || err != nil {
		return nil, err
	}

	fmt.Printf("signature %s\n", s.String())
	//idBytes := 1
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
	// Write ID.
	err := binary.Write(buf, byteOrder, s.idBytes())
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.idBytes())
	// Write versions.
	err = binary.Write(buf, byteOrder, s.majorVersion)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.majorVersion)
	err = binary.Write(buf, byteOrder, s.minorVersion)
	if err != nil {
		return []byte{}, err
	}
	byteCount += binary.Size(s.minorVersion)
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
