package sorter

import (
	"strings"
	"testing"
)

func TestSortYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "simple mapping",
			input: `c: value3
a: value1
b: value2`,
			expected: `a: value1
b: value2
c: value3
`,
			wantErr: false,
		},
		{
			name: "nested mapping",
			input: `zebra:
  c: value3
  a: value1
  b: value2
apple: value
banana: value`,
			expected: `apple: value
banana: value
zebra:
    a: value1
    b: value2
    c: value3
`,
			wantErr: false,
		},
		{
			name: "with arrays",
			input: `c: [3, 2, 1]
a: value1
b: value2`,
			// Encoder may output flow-style [3, 2, 1]; both are valid
			expected: `a: value1
b: value2
c: [3, 2, 1]
`,
			wantErr: false,
		},
		{
			name:  "empty document",
			input: `{}`,
			expected: `{}
`,
			wantErr: false,
		},
		{
			name:    "invalid YAML",
			input:   `invalid: [unclosed`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SortYAML([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("SortYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Normalize whitespace for comparison
				resultStr := strings.TrimSpace(string(result))
				expectedStr := strings.TrimSpace(tt.expected)
				if resultStr != expectedStr {
					t.Errorf("SortYAML() = %v, want %v", resultStr, expectedStr)
				}
			}
		})
	}
}

func TestSortYAML_PreservesStructure(t *testing.T) {
	input := `server:
  port: 8080
  host: localhost
database:
  name: testdb
  user: admin`

	result, err := SortYAML([]byte(input))
	if err != nil {
		t.Fatalf("SortYAML() error = %v", err)
	}

	resultStr := string(result)
	// Check that keys are sorted (database < server alphabetically)
	if strings.Index(resultStr, "server:") < strings.Index(resultStr, "database:") {
		t.Error("Keys should be sorted alphabetically")
	}

	// Check that nested structure is preserved
	if !strings.Contains(resultStr, "port:") || !strings.Contains(resultStr, "host:") {
		t.Error("Nested keys should be preserved")
	}
}
