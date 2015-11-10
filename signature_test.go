package zgok

import (
	"testing"
)

func TestDumpRestore(t *testing.T) {
	// Create signature for testing.
	orig := NewSignature()
	orig.exeSize = 12345678909876
	orig.zipSize = 98765432101234
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
	// Compare data.
	if orig.app != copy.app {
		t.Errorf("Compare app: expected [%v] got [%v]", orig.app, copy.app)
	}
	if orig.major != copy.major {
		t.Errorf("Compare major version: expected [%v] got [%v]",
			orig.major, copy.major)
	}
	if orig.minor != copy.minor {
		t.Errorf("Compare minor version: expected [%v] got [%v]",
			orig.minor, copy.minor)
	}
	if orig.rev != copy.rev {
		t.Errorf("Compare revision: expected [%v] got [%v]",
			orig.rev, copy.rev)
	}
	if orig.exeSize != copy.exeSize {
		t.Errorf("Compare exe size: expected [%v] got [%v]",
			orig.exeSize, copy.exeSize)
	}
	if orig.zipSize != copy.zipSize {
		t.Errorf("Compare zip size: expected [%v] got [%v]",
			orig.zipSize, copy.zipSize)
	}
	// Compare strings.
	if orig.String() != copy.String() {
		t.Errorf("Compare strings: expected [%v] got [%v]",
			orig.String(), copy.String())
	}
}
