// godocall runs go doc -all on all sub folder containing .go files from current path
// The purpose is to update a README.md with all the project documentation for easy access an search on first page.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO Current error is that go doc will not run och gives doc: no buildable Go source
// files in C:\Users\David Nilsson\Dropbox\Dev\Privat\go\src\training\codeChall\leetcode
// like it running in the wrong catalog.

// Flags global struct for argument flags
type flags struct {
	verbose *bool
	path    *string
}

var f flags

func main() {
	f.verbose = flag.Bool("verbose", false, "Enable verbose output")
	f.path = flag.String("path", "", "Root path to run godocall")
	root := "."
	foldermap := make(map[string]struct{})
	const filePostfix = ".go"

	// var b bytes.Buffer
	// b.Grow(32 << (10 * 2))

	flag.Parse()

	if *f.path != "" {
		_, err := os.Stat(*f.path)
		if os.IsNotExist(err) {
			fmt.Printf("Error: The specified path '%s' does not exist.\n", *f.path)
			os.Exit(1)
		} else if err != nil {
			fmt.Printf("Error: Failed to check the specified path '%s': %v\n", *f.path, err)
			os.Exit(1)
		}
		root = *f.path
	}

	if *f.verbose {
		fmt.Println("Starting to travers path:", *f.path)
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == filePostfix {
			if _, ok := foldermap[filepath.Dir(path)]; !ok {
				if *f.verbose {
					fmt.Println("Found files ending with :", filePostfix, "in:", filepath.Dir(path))
				}
				foldermap[filepath.Dir(path)] = struct{}{}
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error: traversing directory: %s\n", err)
		os.Exit(1)
	}

	// buf, err := debug(foldermap)
	buf, err := goDocCatalog(foldermap)
	if err != nil {
		fmt.Printf("Error: scanning directory: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(buf.String())
}

// TODO argument for what command to run
func goDocCatalog(foldermap map[string]struct{}) (out bytes.Buffer, err error) {
	out.Grow(32 << (10 * 2))
	for path := range foldermap {
		if *f.verbose {
			fmt.Printf("Running 'go doc -all' for directory: %s\n", path)
		}

		cmd := exec.Command("go", "doc", "-all", path)
		cmdout, err := cmd.CombinedOutput()
		if err != nil && *f.verbose {
			fmt.Printf("Error running 'go doc -all': %s %s\n", path, err)
		}

		if len(cmdout) > 2 {
			out.Write(cmdout)
		}
	}
	return out, nil
}
