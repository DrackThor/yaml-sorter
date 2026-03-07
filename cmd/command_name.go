package cmd

import (
	"path/filepath"
	"strings"
)

const (
	defaultCommandName = "ysort"
	shortCommandName   = "ysort"
	legacyCommandName  = "yaml-sort"
)

func commandNameFromArg0(arg0 string) string {
	if arg0 == "" {
		return defaultCommandName
	}

	normalized := strings.ReplaceAll(arg0, "\\", "/")
	base := filepath.Base(normalized)
	base = strings.TrimSuffix(base, filepath.Ext(base))

	switch base {
	case shortCommandName:
		return shortCommandName
	case legacyCommandName:
		return shortCommandName
	default:
		return defaultCommandName
	}
}
