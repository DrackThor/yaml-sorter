package version

import (
	"os/exec"
	"runtime/debug"
	"strings"
)

// BuildVersion is injected at build time for release/tag builds.
// Example:
//
//	-ldflags "-X github.com/drackthor/ysort/internal/version.BuildVersion=v1.2.3"
var BuildVersion string

// String returns the current ysort version.
//
// Priority:
// 1. BuildVersion (tag pipeline builds via ldflags)
// 2. Local fallback: <latest-tag>-local.<short-commit>
// 3. Final fallback: v0.0.0-local.unknown
func String() string {
	if normalized := strings.TrimSpace(BuildVersion); normalized != "" {
		return normalized
	}

	latestTag := gitOutput("describe", "--tags", "--abbrev=0", "--match", "v[0-9]*")
	commit := gitOutput("rev-parse", "--short=7", "HEAD")
	if commit == "" {
		commit = shortBuildInfoCommit()
	}

	return formatLocalVersion(latestTag, commit)
}

func gitOutput(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func shortBuildInfoCommit() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok || bi == nil {
		return ""
	}
	for _, setting := range bi.Settings {
		if setting.Key != "vcs.revision" {
			continue
		}
		revision := strings.TrimSpace(setting.Value)
		if revision == "" {
			return ""
		}
		if len(revision) > 7 {
			return revision[:7]
		}
		return revision
	}
	return ""
}

func formatLocalVersion(tag, commit string) string {
	normalizedTag := strings.TrimSpace(tag)
	if normalizedTag == "" {
		normalizedTag = "v0.0.0"
	}

	normalizedCommit := strings.TrimSpace(commit)
	if normalizedCommit == "" {
		normalizedCommit = "unknown"
	}

	return normalizedTag + "-local." + normalizedCommit
}
