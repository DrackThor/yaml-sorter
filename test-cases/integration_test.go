package testcases_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drackthor/yaml-sort/internal/sorter"
	"gopkg.in/yaml.v3"
)

// TestRealWorldInputs runs the sorter on all YAML files in test-cases/inputs/
// and verifies: no error, valid YAML output, and round-trip equality.
func TestRealWorldInputs(t *testing.T) {
	inputsDir := "inputs"
	entries, err := os.ReadDir(inputsDir)
	if err != nil {
		t.Fatalf("reading inputs dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		name := e.Name()
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(inputsDir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read file: %v", err)
			}

			sorted, err := sorter.SortYAML(data)
			if err != nil {
				t.Fatalf("SortYAML: %v", err)
			}

			// Output must be valid YAML
			var out yaml.Node
			if err := yaml.Unmarshal(sorted, &out); err != nil {
				t.Fatalf("sorted output is not valid YAML: %v", err)
			}

			// Round-trip: sorting again should be idempotent (same bytes)
			again, err := sorter.SortYAML(sorted)
			if err != nil {
				t.Fatalf("second SortYAML: %v", err)
			}
			if string(again) != string(sorted) {
				t.Error("sort is not idempotent: second sort produced different output")
			}

			// If expected file exists, diff against it
			expectedPath := filepath.Join("expected", name)
			expected, err := os.ReadFile(expectedPath)
			if err == nil {
				if string(expected) != string(sorted) {
					t.Errorf("output does not match expected file %s", expectedPath)
				}
			}
		})
	}
}
