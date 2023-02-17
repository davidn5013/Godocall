// godocall runs go doc -all on all sub folder containing .go files from current path
// The purpose is to update a README.md with all the project documentation for easy access an search on first page.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	path := flag.String("path", "", "Root path to run godocall")
	flag.Parse()
	root := "."

	if *path != "" {
		_, err := os.Stat(*path)
		if os.IsNotExist(err) {
			fmt.Printf("Error: The specified path '%s' does not exist.\n", *path)
			os.Exit(1)
		} else if err != nil {
			fmt.Printf("Error: Failed to check the specified path '%s': %v\n", *path, err)
			os.Exit(1)
		}
		root = *path
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {

			// TODO argument for non verbose
			if *verbose {
				fmt.Printf("Running 'go doc -all' for directory: %s\n", path)
			}

			// TODO argument for what command to run
			cmd := exec.Command("go", "doc", "-all")

			cmd.Dir = filepath.Dir(path)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running 'go doc -all': %s\n", err)
				return nil
			}

			// NOT I thing this can be much fast i print to large buffer instead
			// I did do this in my Rule 110 source code if remember correct.
			fmt.Println(string(out))

		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error traversing directory: %s\n", err)
		os.Exit(1)
	}
}
