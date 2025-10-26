package cmd

import (
	"fmt"

	"github.com/aykay76/ici/internal/parser"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [workflow-file]",
	Short: "Validate a GitHub Actions workflow",
	Long: `Validate a GitHub Actions workflow file for syntax and semantic errors.

Examples:
  ici validate .github/workflows/test.yml
  ici validate workflow.yml --strict`,
	Args: cobra.ExactArgs(1),
	RunE: validateWorkflow,
}

var strict bool

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&strict, "strict", false, "enable strict validation")
}

func validateWorkflow(cmd *cobra.Command, args []string) error {
	workflowFile := args[0]
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Printf("Validating workflow: %s\n", workflowFile)
	}

	workflow, err := parser.ParseWorkflow(workflowFile)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Basic validation
	if workflow.Name == "" {
		fmt.Println("⚠️  Warning: Workflow has no name")
	}

	if len(workflow.Jobs) == 0 {
		return fmt.Errorf("validation failed: no jobs defined")
	}

	if verbose {
		fmt.Printf("✓ Workflow name: %s\n", workflow.Name)
		fmt.Printf("✓ Jobs: %d\n", len(workflow.Jobs))
		for jobID := range workflow.Jobs {
			fmt.Printf("  - %s\n", jobID)
		}
	}

	fmt.Println("✓ Workflow is valid")
	return nil
}
