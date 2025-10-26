package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Workflow represents a GitHub Actions workflow
type Workflow struct {
	Name string            `yaml:"name"`
	On   interface{}       `yaml:"on"` // Can be string, array, or map
	Jobs map[string]Job    `yaml:"jobs"`
	Env  map[string]string `yaml:"env,omitempty"`
}

// Job represents a single job in a workflow
type Job struct {
	Name    string            `yaml:"name,omitempty"`
	RunsOn  interface{}       `yaml:"runs-on"` // Can be string or array
	Steps   []Step            `yaml:"steps"`
	Env     map[string]string `yaml:"env,omitempty"`
	Needs   interface{}       `yaml:"needs,omitempty"` // Can be string or array
	If      string            `yaml:"if,omitempty"`
	Timeout int               `yaml:"timeout-minutes,omitempty"`
}

// Step represents a single step in a job
type Step struct {
	Name string            `yaml:"name,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Run  string            `yaml:"run,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Env  map[string]string `yaml:"env,omitempty"`
	If   string            `yaml:"if,omitempty"`
}

// ParseWorkflow reads and parses a GitHub Actions workflow file
func ParseWorkflow(filePath string) (*Workflow, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &workflow, nil
}

// GetRunsOn extracts the runs-on value as a string
func (j *Job) GetRunsOn() string {
	switch v := j.RunsOn.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return "ubuntu-latest" // Default
}

// GetNeeds returns job dependencies as a slice
func (j *Job) GetNeeds() []string {
	switch v := j.Needs.(type) {
	case string:
		return []string{v}
	case []interface{}:
		needs := make([]string, 0, len(v))
		for _, n := range v {
			if s, ok := n.(string); ok {
				needs = append(needs, s)
			}
		}
		return needs
	}
	return nil
}
