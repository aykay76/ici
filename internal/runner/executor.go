package runner

import (
	"fmt"

	"github.com/aykay76/ici/internal/container"
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

	// Create container based on runs-on
	mgr := container.NewManager(e.verbose)
	image, err := mgr.MapRunsOn(job.GetRunsOn())
	if err != nil {
		return fmt.Errorf("failed to map runs-on for job %s: %w", jobID, err)
	}

	// Build a simple ContainerConfig: pass job-level env into the container.
	cfg := &container.ContainerConfig{}
	if len(job.Env) > 0 {
		envs := make([]string, 0, len(job.Env))
		for k, v := range job.Env {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
		cfg.Env = envs
	}

	containerID, err := mgr.CreateContainerWithConfig(image, jobID, cfg)
	if err != nil {
		return fmt.Errorf("failed to create container for job %s: %w", jobID, err)
	}
	// Ensure cleanup: stop then remove the container explicitly so lifecycle is clear.
	defer func() {
		// best-effort stop; ignore error to prefer removal
		_ = mgr.StopContainer(containerID)
		_ = mgr.RemoveContainer(containerID)
	}()

	// Execute each step inside the container
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

		if step.Run != "" {
			if err := mgr.RunCommand(containerID, step.Run); err != nil {
				return fmt.Errorf("step %d failed: %w", i+1, err)
			}
		}
	}

	fmt.Printf("✓ Job '%s' completed successfully\n", jobID)
	return nil
}
