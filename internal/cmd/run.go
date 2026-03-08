package cmd

import (
	"fmt"

	"github.com/aykay76/ici/internal/parser"
	"github.com/aykay76/ici/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [workflow-file]",
	Short: "Run a GitHub Actions workflow locally",
	Long: `Execute a GitHub Actions workflow file in local Podman containers.

Examples:
  ici run .github/workflows/test.yml
  ici run .github/workflows/build.yml --job build
  ici run workflow.yml --event push`,
	Args: cobra.ExactArgs(1),
	RunE: runWorkflow,
}

var (
	jobName     string
	eventName   string
	dryRun      bool
	secretsFile string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&jobName, "job", "j", "", "specific job to run (default: all jobs)")
	runCmd.Flags().StringVarP(&eventName, "event", "e", "push", "event that triggers the workflow")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "parse and plan without executing")
	runCmd.Flags().StringVar(&secretsFile, "secrets", "", "path to secrets file (uses ~/.ici/secrets.json by default)")
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	workflowFile := args[0]
	verbose, _ := cmd.Flags().GetBool("verbose")
	secretsFile, _ := cmd.Flags().GetString("secrets")

	if verbose {
		fmt.Printf("Running workflow: %s\n", workflowFile)
		fmt.Printf("Event: %s\n", eventName)
		if jobName != "" {
			fmt.Printf("Job: %s\n", jobName)
		}
		if secretsFile != "" {
			fmt.Printf("Secrets file: %s\n", secretsFile)
		}
	}

	// Parse the workflow
	workflow, err := parser.ParseWorkflow(workflowFile)
	if err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	if verbose {
		fmt.Printf("Parsed workflow: %s\n", workflow.Name)
		fmt.Printf("Jobs found: %d\n", len(workflow.Jobs))
	}

	if dryRun {
		fmt.Println("Dry run mode - workflow parsed successfully")
		return nil
	}

	// Execute the workflow
	var executor *runner.Executor
	if secretsFile != "" {
		executor = runner.NewExecutorWithSecrets(verbose, secretsFile)
	} else {
		executor = runner.NewExecutor(verbose)
	}
	return executor.Run(workflow, jobName, eventName)
}
