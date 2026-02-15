package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// File holds the yaml-sort configuration (e.g. from .yaml-sort.yaml).
type File struct {
	// ListSortKeys defines how to sort lists of objects: for each path (e.g. "spec.egress"),
	// sort the list by the given key (e.g. "name") within each element.
	ListSortKeys []ListSortRule `yaml:"listSortKeys"`
}

// ListSortRule defines a single rule: sort the list at path by each element's key.
type ListSortRule struct {
	Path string `yaml:"path"` // Dot-separated path from root, e.g. "spec.egress"
	Key  string `yaml:"key"`  // Key inside each list element to sort by, e.g. "name"
}

// Load reads a config file from path. Returns nil if the file does not exist.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &f, nil
}
