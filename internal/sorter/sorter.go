package sorter

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// SortYAML sorts a YAML file alphabetically by keys
func SortYAML(data []byte) ([]byte, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML document")
	}

	// Sort the root node
	sortNode(&node.Content[0])

	// Marshal back to YAML
	result, err := yaml.Marshal(&node)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return result, nil
}

// sortNode recursively sorts YAML nodes by their keys
func sortNode(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		// Sort mapping nodes by key
		sortMappingNode(node)
	case yaml.SequenceNode:
		// Recursively sort sequence items
		for _, child := range node.Content {
			sortNode(child)
		}
	}
}

// sortMappingNode sorts a mapping node's key-value pairs by key
func sortMappingNode(node *yaml.Node) {
	if node.Kind != yaml.MappingNode || len(node.Content)%2 != 0 {
		return
	}

	// Extract key-value pairs
	type kvPair struct {
		key   *yaml.Node
		value *yaml.Node
	}

	pairs := make([]kvPair, 0, len(node.Content)/2)
	for i := 0; i < len(node.Content); i += 2 {
		pairs = append(pairs, kvPair{
			key:   node.Content[i],
			value: node.Content[i+1],
		})
		// Recursively sort the value
		sortNode(node.Content[i+1])
	}

	// Sort pairs by key value
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key.Value < pairs[j].key.Value
	})

	// Rebuild Content array with sorted pairs
	node.Content = make([]*yaml.Node, 0, len(pairs)*2)
	for _, pair := range pairs {
		node.Content = append(node.Content, pair.key, pair.value)
	}
}
