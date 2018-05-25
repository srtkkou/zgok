package zgok

import (
	"os"
	"os/exec"
	fpath "path/filepath"
	"runtime"
	"testing"
)

const (
	REAL_EXE_SRC_PATH = "testdata/hello.go"
	REAL_EXE_PATH     = "hello.out"
	DUMMY_EXE_PATH    = "testdata/executable"
)

var (
	exePath string = DUMMY_EXE_PATH
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	defer os.Exit(code)
	teardown()
}

func setup() {
	// Erase unused *.out files.
	paths, _ := fpath.Glob("*.out")
	for _, path := range paths {
		os.Remove(path)
	}
	// Compile testdata/hello.go
	cmd := exec.Command("go", "build", "-o", REAL_EXE_PATH, REAL_EXE_SRC_PATH)
	err := cmd.Run()
	if err == nil {
		exePath = REAL_EXE_PATH
	}
}

func teardown() {
	// Erase unused *.out files.
	paths, _ := fpath.Glob("*.out")
	for _, path := range paths {
		os.Remove(path)
	}
}

func TestBuildRestore(t *testing.T) {
	// Build zgok file.
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
	// Check if exe is not broken.
	if exePath == REAL_EXE_PATH {
		// Get parent path.
		_, filename, _, _ := runtime.Caller(0)
		parentPath, err := fpath.Abs(fpath.Join(filename, ".."))
		// Get expected output.
		exBytes, err := exec.Command(fpath.Join(parentPath, exePath)).Output()
		if err != nil {
			t.Errorf("Failed to run test executable:[%v]", err.Error())
		}
		expected := string(exBytes[:])
		// Get output of the built executable file.
		outBytes, err := exec.Command(fpath.Join(parentPath, outPath)).Output()
		if err != nil {
			t.Errorf("Failed to run built executable:[%v]", err.Error())
		}
		output := string(outBytes[:])
		// Compare the results.
		if expected != output {
			t.Errorf("Different output for built executable:expected [%v] got [%v].",
				expected, output)
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
