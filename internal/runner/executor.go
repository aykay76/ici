package runner

import (
	"fmt"

	"github.com/aykay76/ici/internal/parser"
)

// Executor handles workflow execution
type Executor struct {
	verbose bool
}

// NewExecutor creates a new workflow executor
func NewExecutor(verbose bool) *Executor {
	return &Executor{
		verbose: verbose,
	}
}

// Run executes a workflow
func (e *Executor) Run(workflow *parser.Workflow, jobName string, eventName string) error {
	if e.verbose {
		fmt.Printf("Executing workflow: %s\n", workflow.Name)
		fmt.Printf("Event: %s\n", eventName)
	}

	// If specific job requested, run only that job
	if jobName != "" {
		job, exists := workflow.Jobs[jobName]
		if !exists {
			return fmt.Errorf("job '%s' not found in workflow", jobName)
		}
		return e.runJob(jobName, job)
	}

	// Otherwise run all jobs (TODO: handle dependencies)
	for jobID, job := range workflow.Jobs {
		if err := e.runJob(jobID, job); err != nil {
			return fmt.Errorf("job '%s' failed: %w", jobID, err)
		}
	}

	return nil
}

func (e *Executor) runJob(jobID string, job parser.Job) error {
	if e.verbose {
		fmt.Printf("\n=== Running job: %s ===\n", jobID)
		fmt.Printf("Runs-on: %s\n", job.GetRunsOn())
		fmt.Printf("Steps: %d\n", len(job.Steps))
	}

	// TODO: Create container based on runs-on
	// TODO: Execute each step

	for i, step := range job.Steps {
		if e.verbose {
			fmt.Printf("\nStep %d: %s\n", i+1, step.Name)
			if step.Uses != "" {
				fmt.Printf("  Uses: %s\n", step.Uses)
			}
			if step.Run != "" {
				fmt.Printf("  Run: %s\n", step.Run)
			}
		}
	}

	fmt.Printf("✓ Job '%s' completed successfully\n", jobID)
	return nil
}
