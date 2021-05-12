package main

import (
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

func main() {
	// Parse flags
	flag.Usage = usage
	isHelp := flag.Bool("h", false, "Print usage.")
	isVersion := flag.Bool("v", false, "Print version.")
	flag.Parse()
	// Show usage for [-h]
	if *isHelp {
		usage()
		os.Exit(NORMAL_CODE)
	}
	// Show version for [-v].
	if *isVersion {
		fmt.Println(zgok.Version())
		os.Exit(NORMAL_CODE)
	}
	// Check argument length.
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(ERROR_CODE)
	}
	// Switch by command.
	switch args[0] {
	case "build":
		runBuildCommand(args[1:])
	case "show":
		runShowCommand(args[1:])
	default:
		usage()
		os.Exit(ERROR_CODE)
	}
}

// Usage.
func usage() {
	fmt.Println("Usage: zgok [global flags] <command> [flags]")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Println("  build     : Build zgok executable file.")
	fmt.Println("  show      : Show information in zgok executable file.")
	fmt.Println()
	fmt.Println("global flags:")
	fmt.Println("  -h        : Print this help message.")
	fmt.Println("  -v        : Print version.")
	fmt.Println()
	fmt.Println("build command flags:")
	fmt.Println("  -e string : [REQUIRED] Executable file's path.")
	fmt.Println("  -z string : [REQUIRED] Target paths to add to zip.")
	fmt.Println("  -o string : Output file's path.")
	fmt.Println()
	fmt.Println("show command flags:")
	fmt.Println("  -f        : [REQUIRED] Zgok file's path.")
}

// Run build command.
func runBuildCommand(args []string) {
	// Check argument length.
	if len(args) == 0 {
		usage()
		os.Exit(ERROR_CODE)
	}
	// Parse flags
	var (
		exePath  string
		zipPaths strSlice
		outPath  string
	)
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	fs.StringVar(&exePath, "e", "", "Executable file's path.")
	fs.Var(&zipPaths, "z", "ZIP target paths.")
	fs.StringVar(&outPath, "o", "out", "Output file's path.")
	fs.Parse(args)
	// Validate arguments.
	if exePath == "" || len(zipPaths) == 0 || outPath == "" {
		usage()
		os.Exit(ERROR_CODE)
	}
	// Initialize builder.
	builder := zgok.NewZgokBuilder()
	builder.SetOutPath(outPath)
	err := builder.SetExePath(exePath)
	if err != nil {
		panic(err)
	}
	// Add zip paths.
	for _, zipPath := range zipPaths {
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
	fmt.Printf("Exported %s\n", outPath)
}

// Run show command.
func runShowCommand(args []string) {
	// Check argument length.
	if len(args) == 0 {
		usage()
		os.Exit(ERROR_CODE)
	}
	// Parse flags.
	var (
		filePath string
	)
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.StringVar(&filePath, "f", "", "Zgok file's path.")
	fs.Parse(args)
	// Check file path.
	if filePath == "" {
		flag.Usage()
		os.Exit(ERROR_CODE)
	}
	// Restore zgok file system.
	zfs, err := zgok.RestoreFileSystem(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(ERROR_CODE)
	}
	// Show version.
	fmt.Println("Signature:")
	fmt.Println("  " + zfs.String())
	// Show paths.
	fmt.Println()
	fmt.Println("Paths:")
	for _, path := range zfs.Paths() {
		fmt.Println("  " + path)
	}
}
