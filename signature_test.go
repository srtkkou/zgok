package zgok

import (
	"testing"
)

func TestSignatureDumpRestore(t *testing.T) {
	// Create signature for testing.
	orig := NewSignature()
	orig.SetExeSize(12345678909876)
	orig.SetZipSize(87654321090123)
	// Dump signature to bytes.
	bytes, err := orig.Dump()
	if err != nil {
		t.Errorf("Dump() failed: %v", err)
	}
	// Restore signature from bytes.
	copy, err := RestoreSignature(bytes)
	if err != nil {
		t.Errorf("RestoreSignature() failed: %v", err)
	}
	// Compare strings.
	if orig.String() != copy.String() {
		t.Errorf("Compare strings: expected [%v] got [%v]",
			orig.String(), copy.String())
	}
	// Compare data.
	if orig.ExeSize() != copy.ExeSize() {
		t.Errorf("Compare exe size: expected [%v] got [%v]",
			orig.ExeSize(), copy.ExeSize())
	}
	if orig.ZipSize() != copy.ZipSize() {
		t.Errorf("Compare zip size: expected [%v] got [%v]",
			orig.ZipSize(), copy.ZipSize())
	}

}
