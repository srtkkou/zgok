package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/srtkkou/zgok"
	"os"
)

const (
	NORMAL_CODE = 0
	ERROR_CODE  = 255
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
	isHelp    bool
	isVersion bool
	exePath   string
	zipPaths  strSlice
	outPath   string
}

// Set usage.
func (a *Arguments) setUsage() {
	flag.Usage = func() {
		fmt.Println(zgok.Version())
		fmt.Println()
		fmt.Printf("Usage: %s [-h] [-v]\n", zgok.APP)
		fmt.Printf("  -e exePath -z zipPath [-o outPath]\n")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -v : Print version.")
		fmt.Println("  -h : Print this help message.")
		fmt.Println("  -e string : [REQUIRED] Executable file's path.")
		fmt.Println("  -z string : [REQUIRED] Target paths to add to zip.")
		fmt.Println("  -o string : Output file's path.")
		fmt.Println()
		fmt.Println("Note:")
		fmt.Println("  You can set multiple [-z] arguments.")
	}
}

// Set flags.
func (a *Arguments) setFlags() {
	flag.BoolVar(&a.isHelp, "h", false, "Print usage.")
	flag.BoolVar(&a.isVersion, "v", false, "Print this help message.")
	flag.StringVar(&a.exePath, "e", "", "Executable file's path. *REQUIRED")
	flag.Var(&a.zipPaths, "z", "ZIP target paths. *REQUIRED")
	flag.StringVar(&a.outPath, "o", "out", "Output file's path.")
}

// Verify arguments.
func (a *Arguments) verify() error {
	// Return if [-h] is specified.
	if a.isHelp {
		return nil
	}
	// Return if [--version] is specified.
	if a.isVersion {
		return nil
	}
	// Check [-e] and [-z].
	if a.exePath == "" {
		return errors.New("-e option is required.")
	}
	if len(a.zipPaths) <= 0 {
		return errors.New("-z option is required.")
	}
	return nil
}

func main() {
	// Parse flags
	args := &Arguments{}
	args.setUsage()
	args.setFlags()
	flag.Parse()
	// Verify arguments.
	err := args.verify()
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		flag.Usage()
		os.Exit(ERROR_CODE)
	}
	// Show usage for [-h]
	if args.isHelp {
		flag.Usage()
		os.Exit(NORMAL_CODE)
	}
	// Show version for [--version].
	if args.isVersion {
		fmt.Println(zgok.Version())
		os.Exit(NORMAL_CODE)
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
