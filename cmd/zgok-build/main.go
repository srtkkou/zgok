package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/srtkkou/zgok"
)

// Multiple arguments.
type strSlice []string

// Multiple arguments to string.
func (s *strSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// Add value on multiple argument.
func (s *strSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	// Flag variants.
	var (
		isDebug    bool
		isVerbose bool
		exePath    string
		zipPaths strSlice
		outPath    string
	)
	// Customize usage.
	flag.Usage = func() {
		fmt.Printf("%s version %d.%d.%d\n",
			zgok.APP, zgok.MAJOR, zgok.MINOR, zgok.REV)
		fmt.Printf("Usage: %s -e exePath -z zipPath -o outPath\n",
			zgok.APP)
		flag.PrintDefaults()
	}
	// Parse flags
	flag.BoolVar(&isDebug, "debug", false, "Debug flag.")
	flag.BoolVar(&isVerbose, "verbose", false, "Verbose flag.")
	flag.StringVar(&exePath, "e", "", "Executable file's path. *REQUIRED")
	flag.Var(&zipPaths, "z", "ZIP target paths. *REQUIRED")
	flag.StringVar(&outPath, "o", "out", "Output file's path.")
	flag.Parse()
	// Check arguments.
	if exePath == "" {
		fmt.Fprintln(os.Stderr, "-e option is required.")
		os.Exit(255)
	}
	if len(zipPaths) <= 0 {
		fmt.Fprintln(os.Stderr, "-z option is required.")
		os.Exit(255)
	}
	// Get exe file bytes.
	exeBytes, err := ioutil.ReadFile(exePath)
	if err != nil {
		panic(err)
	}
	// Create zip and get its bytes.
	zipBytes, err := createZipAndGetBytes(zipPaths)
	if err != nil {
		panic(err)
	}
	// Create signature.
	signature := zgok.NewSignature()
	signature.SetExeSize(int64(len(exeBytes)))
	signature.SetZipSize(int64(len(zipBytes)))
	sigBytes, err := signature.Dump()
	if err != nil {
		panic(err)
	}
	// Create out file.
	file, err := os.Open(outPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Append bytes.
	_, err = file.Write(exeBytes)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(zipBytes)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(sigBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Exported %s\n", outPath)
}

// Create zip and get its bytes.
func createZipAndGetBytes(zipPaths strSlice) ([]byte, error) {
	var err error
	// Create new zipper.
	zipper := zgok.NewZipper()
	// Add targets to zip.
	for _, zipPath := range zipPaths {
		err = zipper.Add(zipPath)
		if err != nil {
			zipper.Close()
			return []byte{}, nil
		}
	}
	// Close zip.
	err = zipper.Close()
	if err != nil {
		return []byte{}, err
	}
	// Get zip bytes.
	zipBytes, err := zipper.Bytes()
	if err != nil {
		return []byte{}, err
	}
	return zipBytes, nil
}
