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

func TestSortYAMLK8s(t *testing.T) {
	input := `spec:
  z: 1
  a: 2
metadata:
  name: foo
kind: ConfigMap
apiVersion: v1
status:
  x: 3`

	result, err := SortYAMLK8s([]byte(input))
	if err != nil {
		t.Fatalf("SortYAMLK8s() error = %v", err)
	}
	resultStr := string(result)
	// Root must be in K8s order: apiVersion, kind, metadata, spec, status
	apiVersionPos := strings.Index(resultStr, "apiVersion:")
	kindPos := strings.Index(resultStr, "kind:")
	metadataPos := strings.Index(resultStr, "metadata:")
	specPos := strings.Index(resultStr, "spec:")
	statusPos := strings.Index(resultStr, "status:")
	if apiVersionPos == -1 || kindPos == -1 || metadataPos == -1 || specPos == -1 || statusPos == -1 {
		t.Fatalf("missing one of apiVersion/kind/metadata/spec/status")
	}
	if apiVersionPos >= kindPos || kindPos >= metadataPos || metadataPos >= specPos || specPos >= statusPos {
		t.Errorf("root keys not in K8s order: apiVersion=%d kind=%d metadata=%d spec=%d status=%d",
			apiVersionPos, kindPos, metadataPos, specPos, statusPos)
	}
	// Nested keys under spec should be alphabetical (a before z)
	specSection := resultStr[specPos:]
	if strings.Index(specSection, "a:") > strings.Index(specSection, "z:") {
		t.Error("under spec, keys should be sorted alphabetically")
	}
}

func TestSortYAMLWithOptions_ListSortKeys(t *testing.T) {
	input := `spec:
  egress:
    - name: nv.consul-server.consul-egress-1
      action: allow
    - name: nv.consul-server.consul-egress-0
      action: allow
  ingress:
    - name: nv.ui-ingress-1
      action: allow
    - name: nv.ui-ingress-0
      action: allow
`
	opts := Options{
		ListSortKeys: map[string]string{
			"spec.egress":  "name",
			"spec.ingress": "name",
		},
	}
	result, err := SortYAMLWithOptions([]byte(input), opts)
	if err != nil {
		t.Fatalf("SortYAMLWithOptions() error = %v", err)
	}
	resultStr := string(result)
	// spec.egress should be sorted by name: -egress-0 before -egress-1
	if strings.Index(resultStr, "egress-0") > strings.Index(resultStr, "egress-1") {
		t.Error("spec.egress list should be sorted by name (egress-0 before egress-1)")
	}
	// spec.ingress should be sorted by name: -ingress-0 before -ingress-1
	if strings.Index(resultStr, "ingress-0") > strings.Index(resultStr, "ingress-1") {
		t.Error("spec.ingress list should be sorted by name (ingress-0 before ingress-1)")
	}
}
