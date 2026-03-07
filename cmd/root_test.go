package cmd

import "testing"

func TestCommandNameFromArg0(t *testing.T) {
	tests := []struct {
		name     string
		arg0     string
		expected string
	}{
		{
			name:     "yaml-sort binary",
			arg0:     "/usr/local/bin/yaml-sort",
			expected: shortCommandName,
		},
		{
			name:     "ysort binary",
			arg0:     "/usr/local/bin/ysort",
			expected: shortCommandName,
		},
		{
			name:     "ysort windows style path",
			arg0:     `C:\Tools\ysort.exe`,
			expected: shortCommandName,
		},
		{
			name:     "unknown binary falls back to default",
			arg0:     "/usr/local/bin/custom-name",
			expected: defaultCommandName,
		},
		{
			name:     "empty arg falls back to default",
			arg0:     "",
			expected: defaultCommandName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := commandNameFromArg0(tc.arg0)
			if got != tc.expected {
				t.Fatalf("commandNameFromArg0(%q) = %q, want %q", tc.arg0, got, tc.expected)
			}
		})
	}
}
