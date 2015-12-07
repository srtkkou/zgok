package zgok

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	defer os.Exit(code)
	teardown()
}

func teardown() {
	paths, _ := filepath.Glob("*.out")
	for _, path := range paths {
		os.Remove(path)
	}
}

func TestBuildRestore(t *testing.T) {
	// Build zgok file.
	exePath := "testdata/executable"
	outPath := "builder_test.out"
	builder := NewZgokBuilder()
	err := builder.SetExePath(exePath)
	if err != nil {
		t.Errorf("SetExePath():error=[%v]", err)
	}
	err = builder.AddZipPath("testdata/foo")
	if err != nil {
		t.Errorf("AddZipPath():error=[%v]", err)
	}
	err = builder.AddZipPath("testdata/dir")
	if err != nil {
		t.Errorf("AddZipPath():error=[%v]", err)
	}
	builder.SetOutPath(outPath)
	err = builder.Build()
	if err != nil {
		t.Errorf("Build():error=[%v]", err)
	}
	// Load zgok filesystem.
	zfs, err := RestoreFileSystem(outPath)
	if err != nil {
		t.Errorf("RestoreFileSystem():error=[%v]", err)
	}
	// Compare exe size.
	exeStat, err := os.Stat(exePath)
	if err != nil {
		t.Errorf("Failed to read stat of [%v].", exePath)
	}
	if exeStat.Size() != zfs.Signature().ExeSize() {
		t.Errorf("Exe size:expected [%v] got [%v].",
			exeStat.Size(), zfs.Signature().ExeSize())
	}
	// Compare included zip paths.
	zipPaths := []string{
		"testdata/dir/bar",
		"testdata/dir/baz",
		"testdata/foo",
	}
	for i, path := range zfs.Paths() {
		if zipPaths[i] != path {
			t.Errorf("Paths():expected [%v] got [%v].", zipPaths[i], path)
		}
	}
}

func TestBuilderSetInvalidExePath(t *testing.T) {
	builder := NewZgokBuilder()
	err := builder.SetExePath("invalid")
	if err == nil {
		t.Errorf("Expected error on setting invalid exe path.")
	}
}

func TestBuilderAddInvalidZipPath(t *testing.T) {
	builder := NewZgokBuilder()
	err := builder.AddZipPath("invalid")
	if err == nil {
		t.Errorf("Expected error on adding invalid zip path.")
	}
}

func TestBuilderBuildWithoutExePath(t *testing.T) {
	builder := NewZgokBuilder()
	builder.AddZipPath("testdata/foo")
	builder.AddZipPath("testdata/dir")
	builder.SetOutPath("builder_test.out")
	err := builder.Build()
	if err == nil {
		t.Errorf("Expected error on build without exe path.")
	}
}

func TestBuilderBuildWithoutZipPath(t *testing.T) {
	builder := NewZgokBuilder()
	builder.SetExePath("testdata/executable")
	builder.SetOutPath("builder_test.out")
	err := builder.Build()
	if err == nil {
		t.Errorf("Expected error on build without zip path.")
	}
}
