package sorter

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// K8sRootKeyOrder defines the preferred order of top-level keys in a Kubernetes
// manifest. Keys not in this list are sorted alphabetically after these.
var K8sRootKeyOrder = []string{"apiVersion", "kind", "metadata", "spec", "data", "status"}

// Options configures how YAML is sorted.
type Options struct {
	// K8sRoot: root mapping uses fixed K8s key order (apiVersion, kind, metadata, spec, …).
	K8sRoot bool
	// ListSortKeys: for each path (e.g. "spec.egress"), sort that list by the given key (e.g. "name") in each element.
	// Path is dot-separated from document root, e.g. "spec.ingress", "spec.egress".
	ListSortKeys map[string]string // path -> key
}

// SortYAML sorts a YAML document recursively: at each level, mapping keys are
// sorted alphabetically, and we recurse into each value (and into sequence
// elements) so that nested maps and lists are sorted too.
func SortYAML(data []byte) ([]byte, error) {
	return SortYAMLWithOptions(data, Options{})
}

// SortYAMLK8s sorts a YAML document like SortYAML, but the root mapping (top-level
// keys) is ordered for Kubernetes manifests: apiVersion, kind, metadata, spec, …
// Everything under those keys is still sorted alphabetically (recursive).
func SortYAMLK8s(data []byte) ([]byte, error) {
	return SortYAMLWithOptions(data, Options{K8sRoot: true})
}

// SortYAMLWithOptions sorts a YAML document using the given options (K8s root order,
// and optional list sort keys from a config file).
func SortYAMLWithOptions(data []byte, opts Options) ([]byte, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML document")
	}

	root := node.Content[0]
	normalizeLeadingComments(root, data)
	sortNodeWithPath(root, nil, opts)

	result, err := yaml.Marshal(&node)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return result, nil
}

// sortNodeWithPath recursively sorts the tree. path is the dot-separated path from
// document root to this node (e.g. ["spec", "egress"]). Used to apply listSortKeys.
func sortNodeWithPath(node *yaml.Node, path []string, opts Options) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		sortMappingNodeWithPath(node, path, opts)
	case yaml.SequenceNode:
		sortSequenceNodeWithPath(node, path, opts)
	}
}

func sortMappingNodeWithPath(node *yaml.Node, path []string, opts Options) {
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}
	kvPairs := extractKeyValuePairs(node)
	for _, p := range kvPairs {
		sortNodeWithPath(p.value, append(path, p.key.Value), opts)
	}
	// Root mapping and K8s mode: use fixed key order; otherwise alphabetical
	if len(path) == 0 && opts.K8sRoot {
		sort.Slice(kvPairs, func(i, j int) bool {
			return k8sRootKeyLess(kvPairs[i].key.Value, kvPairs[j].key.Value)
		})
	} else {
		sort.Slice(kvPairs, func(i, j int) bool {
			return kvPairs[i].key.Value < kvPairs[j].key.Value
		})
	}
	rebuildMappingContent(node, kvPairs)
}

func sortSequenceNodeWithPath(node *yaml.Node, path []string, opts Options) {
	if node.Kind != yaml.SequenceNode {
		return
	}
	pathStr := strings.Join(path, ".")
	if key, ok := opts.ListSortKeys[pathStr]; ok {
		// Sort this list by each element's key (e.g. "name")
		sort.Slice(node.Content, func(i, j int) bool {
			vi := getScalarFromMapping(node.Content[i], key)
			vj := getScalarFromMapping(node.Content[j], key)
			return vi < vj
		})
	}
	for _, child := range node.Content {
		sortNodeWithPath(child, path, opts)
	}
}

// getScalarFromMapping returns the scalar value for key in the mapping node, or "" if not found.
func getScalarFromMapping(node *yaml.Node, key string) string {
	if node == nil || node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			v := node.Content[i+1]
			if v.Kind == yaml.ScalarNode {
				return v.Value
			}
			return ""
		}
	}
	return ""
}

func k8sRootKeyLess(a, b string) bool {
	idxA := indexOfK8sRootKey(a)
	idxB := indexOfK8sRootKey(b)
	if idxA >= 0 && idxB >= 0 {
		return idxA < idxB
	}
	if idxA >= 0 {
		return true
	}
	if idxB >= 0 {
		return false
	}
	return a < b
}

func indexOfK8sRootKey(key string) int {
	for i, k := range K8sRootKeyOrder {
		if k == key {
			return i
		}
	}
	return -1
}

type kvPair struct {
	key   *yaml.Node
	value *yaml.Node
}

func extractKeyValuePairs(node *yaml.Node) []kvPair {
	n := len(node.Content) / 2
	pairs := make([]kvPair, 0, n)
	for i := 0; i < len(node.Content); i += 2 {
		pairs = append(pairs, kvPair{
			key:   node.Content[i],
			value: node.Content[i+1],
		})
	}
	return pairs
}

func rebuildMappingContent(node *yaml.Node, pairs []kvPair) {
	node.Content = make([]*yaml.Node, 0, len(pairs)*2)
	for _, p := range pairs {
		node.Content = append(node.Content, p.key, p.value)
	}
}

// normalizeLeadingComments ensures comments that appear directly above nodes in
// source text stay attached to those nodes, including blank separators inside a
// comment block.
func normalizeLeadingComments(node *yaml.Node, data []byte) {
	lines := strings.Split(string(data), "\n")
	normalizeNodeLeadingComments(node, lines)
}

func normalizeNodeLeadingComments(node *yaml.Node, lines []string) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		normalizeMappingLeadingComments(node, lines)
		for i := 1; i < len(node.Content); i += 2 {
			normalizeNodeLeadingComments(node.Content[i], lines)
		}
	case yaml.SequenceNode:
		normalizeSequenceLeadingComments(node, lines)
		for _, item := range node.Content {
			normalizeNodeLeadingComments(item, lines)
		}
	case yaml.DocumentNode:
		for _, child := range node.Content {
			normalizeNodeLeadingComments(child, lines)
		}
	}
}

func normalizeMappingLeadingComments(node *yaml.Node, lines []string) {
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}

	headComments := make([]string, 0, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		head := extractLeadingCommentBlock(lines, key.Line)
		headComments = append(headComments, head)
		if head != "" {
			key.HeadComment = head
		}
	}

	for i := 0; i < len(node.Content)-2; i += 2 {
		nextIdx := (i / 2) + 1
		if headComments[nextIdx] == "" {
			continue
		}
		node.Content[i].FootComment = ""
		node.Content[i+1].FootComment = ""
	}
}

func normalizeSequenceLeadingComments(node *yaml.Node, lines []string) {
	if node.Kind != yaml.SequenceNode {
		return
	}

	headComments := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		head := extractLeadingCommentBlock(lines, item.Line)
		headComments = append(headComments, head)
		if head != "" {
			item.HeadComment = head
		}
	}

	for i := 0; i < len(node.Content)-1; i++ {
		if headComments[i+1] == "" {
			continue
		}
		node.Content[i].FootComment = ""
	}
}

func extractLeadingCommentBlock(lines []string, line int) string {
	if line <= 1 || line > len(lines) {
		return ""
	}

	collected := make([]string, 0)
	hasComment := false
	for i := line - 2; i >= 0; i-- {
		current := lines[i]
		trimmed := strings.TrimSpace(current)
		if trimmed == "" {
			collected = append(collected, "")
			continue
		}

		leftTrimmed := strings.TrimLeft(current, " \t")
		if strings.HasPrefix(leftTrimmed, "#") {
			collected = append(collected, leftTrimmed)
			hasComment = true
			continue
		}

		break
	}

	if !hasComment {
		return ""
	}

	// We only keep blank separators that occur after at least one comment line.
	for len(collected) > 0 && collected[len(collected)-1] == "" {
		collected = collected[:len(collected)-1]
	}

	for i, j := 0, len(collected)-1; i < j; i, j = i+1, j-1 {
		collected[i], collected[j] = collected[j], collected[i]
	}

	for i := range collected {
		if collected[i] == "" {
			collected[i] = "#"
		}
	}

	return strings.Join(collected, "\n")
}
