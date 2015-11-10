package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/srtkkou/zgok"
	"os"
)

const (
	ERROR_CODE = 255
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

// Arguments.
type Arguments struct {
	isDebug   bool
	isVerbose bool
	exePath   string
	zipPaths  strSlice
	outPath   string
}

// Verify arguments.
func (a *Arguments) verify() error {
	if a.exePath == "" {
		return errors.New("-e option is required.")
	}
	if len(a.zipPaths) <= 0 {
		return errors.New("-z option is required.")
	}
	return nil
}

func main() {
	// Customize usage.
	flag.Usage = func() {
		fmt.Printf("%s version %d.%d.%d\n",
			zgok.APP, zgok.MAJOR, zgok.MINOR, zgok.REV)
		fmt.Printf("Usage: %s -e exePath -z zipPath -o outPath\n",
			zgok.APP)
		flag.PrintDefaults()
	}
	// Parse flags
	args := &Arguments{}
	flag.BoolVar(&args.isDebug, "debug", false, "Debug flag.")
	flag.BoolVar(&args.isVerbose, "verbose", false, "Verbose flag.")
	flag.StringVar(&args.exePath, "e", "", "Executable file's path. *REQUIRED")
	flag.Var(&args.zipPaths, "z", "ZIP target paths. *REQUIRED")
	flag.StringVar(&args.outPath, "o", "out", "Output file's path.")
	flag.Parse()
	// Verify arguments.
	err := args.verify()
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		flag.Usage()
		os.Exit(ERROR_CODE)
	}
	// Initialize builder.
	builder := zgok.NewZgokBuilder()
	// Set out path.
	builder.SetOutPath(args.outPath)
	// Set exe path.
	err = builder.SetExePath(args.exePath)
	if err != nil {
		panic(err)
	}
	// Add zip paths.
	for _, zipPath := range args.zipPaths {
		err = builder.AddZipPath(zipPath)
		if err != nil {
			panic(err)
		}
	}
	// Build zgok file.
	err = builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Exported %s\n", args.outPath)
}
