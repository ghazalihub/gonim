package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/go2nim/internal/loader"
	"github.com/user/go2nim/internal/translator"
)

func main() {
	outputDir := flag.String("o", ".", "output directory")
	flag.Parse()

	patterns := flag.Args()
	if len(patterns) == 0 {
		fmt.Println("Usage: go2nim [-o output_dir] <package_patterns>")
		os.Exit(1)
	}

	prog, err := loader.Load(patterns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading packages: %v\n", err)
		os.Exit(1)
	}

	nimFiles := translator.Translate(prog)
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	for i, nf := range nimFiles {
		pkg := prog.Packages[i]
		// Use package name for the Nim file
		outPath := filepath.Join(*outputDir, pkg.Name+".nim")
		err := os.WriteFile(outPath, []byte(nf.Render(0)), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
		} else {
			fmt.Printf("Generated %s\n", outPath)
		}
	}
}
