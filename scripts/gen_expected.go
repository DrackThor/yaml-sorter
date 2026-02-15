// gen_expected reads each YAML file in test-cases/inputs/ and writes
// sorted output to test-cases/expected/. Run from repo root:
//
//	go run ./scripts/gen_expected.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drackthor/yaml-sort/internal/sorter"
)

func main() {
	inputsDir := "test-cases/inputs"
	expectedDir := "test-cases/expected"

	if err := os.MkdirAll(expectedDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(inputsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(inputsDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			continue
		}

		sorted, err := sorter.SortYAML(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sort %s: %v\n", path, err)
			continue
		}

		outPath := filepath.Join(expectedDir, e.Name())
		if err := os.WriteFile(outPath, sorted, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
			continue
		}
		fmt.Println(outPath)
	}
}
