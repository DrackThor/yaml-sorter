package sorter

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// K8sRootKeyOrder defines the preferred order of top-level keys in a Kubernetes
// manifest. Keys not in this list are sorted alphabetically after these.
var K8sRootKeyOrder = []string{"apiVersion", "kind", "metadata", "spec", "data", "status"}

// SortYAML sorts a YAML document recursively: at each level, mapping keys are
// sorted alphabetically, and we recurse into each value (and into sequence
// elements) so that nested maps and lists are sorted too.
func SortYAML(data []byte) ([]byte, error) {
	return sortYAML(data, false)
}

// SortYAMLK8s sorts a YAML document like SortYAML, but the root mapping (top-level
// keys) is ordered for Kubernetes manifests: apiVersion, kind, metadata, spec, …
// Everything under those keys is still sorted alphabetically (recursive).
func SortYAMLK8s(data []byte) ([]byte, error) {
	return sortYAML(data, true)
}

// sortYAML unmarshals the YAML into a tree of yaml.Nodes, sorts the tree,
// then marshals it back to bytes. If k8sRoot is true, the root mapping uses
// K8s key order; otherwise all mappings are sorted alphabetically.
func sortYAML(data []byte, k8sRoot bool) ([]byte, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML document")
	}

	root := node.Content[0]
	if k8sRoot && root.Kind == yaml.MappingNode {
		sortMappingNodeK8sRoot(root)
	} else {
		sortNode(root)
	}

	result, err := yaml.Marshal(&node)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return result, nil
}

// sortNode recursively sorts the YAML tree: for MappingNodes we sort by key and
// recurse into each value; for SequenceNodes we recurse into each element.
// Other node kinds (scalars, etc.) are left as-is.
func sortNode(node *yaml.Node) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		sortMappingNode(node)
	case yaml.SequenceNode:
		for _, child := range node.Content {
			sortNode(child)
		}
	}
}

// sortMappingNode sorts a mapping’s key-value pairs alphabetically by key,
// and recursively sorts each value (so nested maps and lists are sorted too).
func sortMappingNode(node *yaml.Node) {
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}
	kvPairs := extractKeyValuePairs(node)
	for _, p := range kvPairs {
		sortNode(p.value)
	}
	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].key.Value < kvPairs[j].key.Value
	})
	rebuildMappingContent(node, kvPairs)
}

// sortMappingNodeK8sRoot sorts only the root mapping: keys follow K8s order
// (apiVersion, kind, metadata, spec, …), then the rest alphabetically.
// Each value is recursively sorted alphabetically via sortNode.
func sortMappingNodeK8sRoot(node *yaml.Node) {
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}
	kvPairs := extractKeyValuePairs(node)
	for _, p := range kvPairs {
		sortNode(p.value)
	}
	sort.Slice(kvPairs, func(i, j int) bool {
		return k8sRootKeyLess(kvPairs[i].key.Value, kvPairs[j].key.Value)
	})
	rebuildMappingContent(node, kvPairs)
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
