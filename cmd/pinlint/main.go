// Command pinlint fails when any workflow uses an action that is not pinned by commit SHA.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/EquinoxWN/secure-supply-chain/internal/pinlint"
)

func main() {
	dir := flag.String("dir", ".github/workflows", "directory holding workflow files")
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "pinlint: unexpected arguments; use -dir <path>")
		os.Exit(2)
	}
	n, err := run(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pinlint:", err)
		os.Exit(2)
	}
	if n > 0 {
		os.Exit(1)
	}
}

// run checks every workflow in dir and returns the number of findings.
func run(dir string) (int, error) {
	var files []string
	for _, pattern := range []string{"*.yml", "*.yaml"} {
		m, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return 0, err
		}
		files = append(files, m...)
	}
	if len(files) == 0 {
		return 0, fmt.Errorf("no workflow files in %s", dir)
	}
	total := 0
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		findings, err := pinlint.Check(f, path)
		f.Close()
		if err != nil {
			return 0, err
		}
		for _, x := range findings {
			fmt.Println(x)
		}
		total += len(findings)
	}
	fmt.Printf("pinlint: %d workflow file(s), %d unpinned reference(s)\n", len(files), total)
	return total, nil
}
