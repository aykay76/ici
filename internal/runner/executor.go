package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/aykay76/ici/internal/container"
	"github.com/aykay76/ici/internal/parser"
	"github.com/aykay76/ici/internal/secrets"
)

// Executor handles workflow execution
type Executor struct {
	verbose     bool
	secretsFile string
}

// NewExecutor creates a new workflow executor
func NewExecutor(verbose bool) *Executor {
	return &Executor{
		verbose:     verbose,
		secretsFile: "", // Use default location
	}
}

// NewExecutorWithSecrets creates a new workflow executor with a custom secrets file location
func NewExecutorWithSecrets(verbose bool, secretsFile string) *Executor {
	return &Executor{
		verbose:     verbose,
		secretsFile: secretsFile,
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
		return e.runJob(workflow, jobName, job)
	}

	// Otherwise run all jobs (TODO: handle dependencies)
	for jobID, job := range workflow.Jobs {
		if err := e.runJob(workflow, jobID, job); err != nil {
			return fmt.Errorf("job '%s' failed: %w", jobID, err)
		}
	}

	return nil
}

func (e *Executor) runJob(workflow *parser.Workflow, jobID string, job parser.Job) error {
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

	// Build a simple ContainerConfig: pass job-level env into the container and mount the current workspace.
	cfg := &container.ContainerConfig{}

	// Load secrets from the local secret store
	var secretStore secrets.Store
	if e.secretsFile != "" {
		secretStore = secrets.NewFileStore(e.secretsFile)
	} else {
		secretStore = secrets.DefaultStore()
	}
	storedSecrets, err := secretStore.GetAll()
	if err != nil {
		// Log but continue - missing secrets file is not fatal
		if e.verbose {
			fmt.Printf("Warning: could not load secrets: %v\n", err)
		}
	}

	// Merge workflow, job env, and secrets into the container config env
	// Order: workflow env -> job env -> step env -> stored secrets (secrets can override)
	envMap := map[string]string{}
	if workflow.Env != nil {
		for k, v := range workflow.Env {
			envMap[k] = v
		}
	}
	if job.Env != nil {
		for k, v := range job.Env {
			envMap[k] = v
		}
	}
	// Add stored secrets (can override env vars)
	for k, v := range storedSecrets {
		envMap[k] = v
	}

	if len(envMap) > 0 {
		envs := make([]string, 0, len(envMap))
		for k, v := range envMap {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
		cfg.Env = envs
		if e.verbose {
			fmt.Printf("Environment variables loaded: %d (workflow) + %d (job) + %d (secrets)\n",
				len(workflow.Env), len(job.Env), len(storedSecrets))
		}
	}

	// Mount the host's current working directory into the container at /workspace
	// Use a simple bind mount: <hostcwd>:/workspace:rw
	if cwd, err := os.Getwd(); err == nil {
		cfg.Volumes = []string{fmt.Sprintf("%s:/workspace:rw", cwd)}
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
			// Use job timeout (minutes) if specified; convert to seconds for the container RunCommandWithTimeout
			timeoutSeconds := 0
			if job.Timeout > 0 {
				timeoutSeconds = job.Timeout * 60
			}

			// Merge env for this step: start from container config env (workflow+job) and overlay step.Env
			execEnv := make([]string, 0)
			if cfg != nil && len(cfg.Env) > 0 {
				execEnv = append(execEnv, cfg.Env...)
			}
			if step.Env != nil {
				for k, v := range step.Env {
					execEnv = append(execEnv, fmt.Sprintf("%s=%s", k, v))
				}
			}

			// Determine working directory inside container; if step.WorkingDirectory is relative, map it under /workspace
			workDir := step.WorkingDirectory
			if workDir == "" {
				// default to workspace root
				workDir = "/workspace"
			} else if !strings.HasPrefix(workDir, "/") {
				// relative path -> under /workspace
				workDir = fmt.Sprintf("/workspace/%s", workDir)
			}

			if err := mgr.RunCommandWithOptions(containerID, step.Run, execEnv, workDir, timeoutSeconds); err != nil {
				return fmt.Errorf("step %d failed: %w", i+1, err)
			}
		}
	}

	fmt.Printf("✓ Job '%s' completed successfully\n", jobID)
	return nil
}
