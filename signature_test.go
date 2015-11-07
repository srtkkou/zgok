package zgok

import (
	"testing"
	//"fmt"
	//"encoding/binary"
	//"bytes"
)

func TestDumpRestore(t *testing.T) {
	// Create signature for testing.
	sig := NewSignature()
	sig.exeSize = 12345678901234
	sig.zipSize = 98765432123456
	// Dump signature to bytes.
	bytes, err := sig.Dump()
	if err != nil {
		t.Errorf("Dump() failed: %v", err)
	}
	// Restore signature from bytes.
	sig2, err := RestoreSignature(bytes)
	if err != nil {
		t.Errorf("RestoreSignature() failed: %v", err)
	}
	if sig.id != sig2.id {
		t.Errorf("Compare id: expected [%v] got [%v]", sig.id, sig2.id)
	}
	t.Errorf("%v %v", bytes, err)
}
