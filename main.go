// Package main godocs run go doc -all och all sub folder containing .go files from current path
/* The purpose is to update a README.md with all the
   project documentation for easy access an search on first package
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// TODO argument for path
	root := "."

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {

			// TODO argument for non verbose
			fmt.Printf("Running 'go doc -all' for directory: %s\n", path)

			// TODO argument for what command to run
			cmd := exec.Command("go", "doc", "-all")

			cmd.Dir = filepath.Dir(path)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running 'go doc -all': %s\n", err)
				return nil
			}

			fmt.Println(string(out))

		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error traversing directory: %s\n", err)
		os.Exit(1)
	}
}
