package cmd

import (
	"fmt"
	"os"

	"github.com/drackthor/yaml-sort/internal/config"
	"github.com/drackthor/yaml-sort/internal/sorter"
	"github.com/spf13/cobra"
)

var (
	inplace    bool
	output     string
	k8sMode    bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "yaml-sort [file]",
	Short: "A tool to sort YAML files",
	Long: `yaml-sort is a CLI tool that sorts YAML files alphabetically
by their keys while preserving the structure and comments where possible.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		inputFile := args[0]

		// Validate flags
		if inplace && output != "" {
			return fmt.Errorf("cannot use both -i and -o flags together")
		}

		// Read input file
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Build sort options (K8s root + optional list sort keys from config)
		opts := sorter.Options{K8sRoot: k8sMode}
		if configPath != "" {
			cfg, loadErr := config.Load(configPath)
			if loadErr != nil {
				return fmt.Errorf("config: %w", loadErr)
			}
			if cfg != nil && len(cfg.ListSortKeys) > 0 {
				opts.ListSortKeys = make(map[string]string, len(cfg.ListSortKeys))
				for _, r := range cfg.ListSortKeys {
					opts.ListSortKeys[r.Path] = r.Key
				}
			}
		}
		sorted, err := sorter.SortYAMLWithOptions(content, opts)
		if err != nil {
			return fmt.Errorf("failed to sort YAML: %w", err)
		}

		// Write output
		if inplace {
			if err := os.WriteFile(inputFile, sorted, 0644); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
			fmt.Printf("Successfully sorted %s in-place\n", inputFile)
		} else if output != "" {
			if err := os.WriteFile(output, sorted, 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("Successfully sorted %s -> %s\n", inputFile, output)
		} else {
			// Write to stdout
			fmt.Print(string(sorted))
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "sort file in-place, replacing the original file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "write sorted output to specified file")
	rootCmd.Flags().BoolVarP(&k8sMode, "k8s", "k", false, "Kubernetes manifest mode: root keys in fixed order (apiVersion, kind, metadata, spec, …), rest alphabetical")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file defining list sort keys (e.g. sort spec.egress by name)")
}
