package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aykay76/ici/internal/parser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var parseCmd = &cobra.Command{
	Use:   "parse [workflow-file]",
	Short: "Parse and display a GitHub Actions workflow",
	Long: `Parse a GitHub Actions workflow file and display its structure.
Useful for validating syntax and understanding workflow structure.

Examples:
  ici parse .github/workflows/test.yml
  ici parse workflow.yml --format json`,
	Args: cobra.ExactArgs(1),
	RunE: parseWorkflow,
}

var outputFormat string

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "yaml", "output format (yaml, json)")
}

func parseWorkflow(cmd *cobra.Command, args []string) error {
	workflowFile := args[0]
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Printf("Parsing workflow: %s\n", workflowFile)
	}

	workflow, err := parser.ParseWorkflow(workflowFile)
	if err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Output the parsed workflow
	var output []byte
	switch outputFormat {
	case "json":
		output, err = json.MarshalIndent(workflow, "", "  ")
	case "yaml":
		output, err = yaml.Marshal(workflow)
	default:
		return fmt.Errorf("unsupported format: %s (use yaml or json)", outputFormat)
	}

	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
