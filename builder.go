package zgok

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Builder interface.
type Builder interface {
	SetExePath(exePath string) error
	AddZipPath(zipPath string) error
	SetOutPath(outPath string)
	Build() error
}

// Zgok builder
type zgokBuilder struct {
	exePath  string   // Executable file's path.
	zipPaths []string // Zip file path.
	outPath  string   // Output file path.
	exeBytes *[]byte  // Bytes of the executable file.
	zipBytes *[]byte  // Bytes of the zip file.
	sigBytes *[]byte  // Bytes of the signature.
}

// Initialize new zgok builder.
func NewZgokBuilder() Builder {
	b := &zgokBuilder{}
	return b
}

// Set executable file path.
func (b *zgokBuilder) SetExePath(exePath string) error {
	_, err := os.Stat(exePath)
	if err != nil {
		return err
	}
	b.exePath = exePath
	return nil
}

// Add paths to add to zip.
func (b *zgokBuilder) AddZipPath(zipPath string) error {
	_, err := os.Stat(zipPath)
	if err != nil {
		return err
	}
	b.zipPaths = append(b.zipPaths, zipPath)
	return nil
}

// Set output path.
func (b *zgokBuilder) SetOutPath(outPath string) {
	b.outPath = outPath
}

// Build zgok file.
func (b *zgokBuilder) Build() error {
	// Set exe file bytes.
	err := b.setExeBytes()
	if err != nil {
		return err
	}
	// Set zip file bytes.
	err = b.setZipBytes()
	if err != nil {
		return err
	}
	// Set signature bytes.
	err = b.setSignatureBytes()
	if err != nil {
		return err
	}
	// Create out file.
	err = b.createOutFile()
	if err != nil {
		return err
	}
	return nil
}

// Set exe file bytes.
func (b *zgokBuilder) setExeBytes() error {
	exeBytes, err := ioutil.ReadFile(b.exePath)
	if err != nil {
		return err
	}
	b.exeBytes = &exeBytes
	return nil
}

// Set zip file bytes.
func (b *zgokBuilder) setZipBytes() error {
	// Check if zip paths are empty.
	if len(b.zipPaths) == 0 {
		return fmt.Errorf("Zip paths not set.")
	}
	var err error
	// Create new zipper.
	zipper := NewZipper()
	// Add targets to zip.
	for _, zipPath := range b.zipPaths {
		err = zipper.Add(zipPath)
		if err != nil {
			zipper.Close()
			return err
		}
	}
	// Close zip.
	err = zipper.Close()
	if err != nil {
		return err
	}
	// Get zip bytes.
	zipBytes, err := zipper.Bytes()
	if err != nil {
		return err
	}
	b.zipBytes = &zipBytes
	return nil
}

// Set signature bytes.
func (b *zgokBuilder) setSignatureBytes() error {
	// Check if exeBytes and zipBytes are set.
	if len(*b.exeBytes) <= 0 {
		return fmt.Errorf("Exe bytes not set.")
	}
	if len(*b.zipBytes) <= 0 {
		return fmt.Errorf("Zip bytes not set.")
	}
	// Get sizes.
	exeSize := int64(len(*b.exeBytes))
	zipSize := int64(len(*b.zipBytes))
	// Create signature.
	signature := NewSignature()
	signature.SetExeSize(exeSize)
	signature.SetZipSize(zipSize)
	// Dump signature.
	sigBytes, err := signature.Dump()
	if err != nil {
		return err
	}
	b.sigBytes = &sigBytes
	return nil
}

// Create out file.
func (b *zgokBuilder) createOutFile() error {
	// Create out file.
	file, err := os.Create(b.outPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Append bytes.
	_, err = file.Write(*b.exeBytes)
	if err != nil {
		return err
	}
	_, err = file.Write(*b.zipBytes)
	if err != nil {
		return err
	}
	_, err = file.Write(*b.sigBytes)
	if err != nil {
		return err
	}
	// Synchronize file.
	file.Sync()
	return nil
}
