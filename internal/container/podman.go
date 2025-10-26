package container

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// execCommand is a package-level variable so tests can override command execution.
var execCommand = exec.Command

// Manager handles Podman container operations
type Manager struct {
	verbose bool
	cli     string // detected CLI: podman or docker
}

// NewManager creates a new container manager
func NewManager(verbose bool) *Manager {
	m := &Manager{
		verbose: verbose,
		cli:     "",
	}
	// detect available container CLI
	if path, err := exec.LookPath("podman"); err == nil {
		m.cli = path
	} else if path, err := exec.LookPath("docker"); err == nil {
		m.cli = path
	}

	return m
}

// MapRunsOn converts GitHub Actions runs-on to container images
func (m *Manager) MapRunsOn(runsOn string) (string, error) {
	// Map GitHub runners to container images
	imageMap := map[string]string{
		"ubuntu-latest": "ubuntu:22.04",
		"ubuntu-22.04":  "ubuntu:22.04",
		"ubuntu-20.04":  "ubuntu:20.04",
		// TODO: Add more mappings
	}

	if image, ok := imageMap[runsOn]; ok {
		return image, nil
	}

	return "", fmt.Errorf("unsupported runs-on: %s", runsOn)
}

// CreateContainer creates a new Podman (or Docker) container and returns its ID.
// It pulls the image, creates the container (keeps it running) and starts it.
func (m *Manager) CreateContainer(image string, name string) (string, error) {
	if m.verbose {
		fmt.Printf("Creating container: %s (image: %s)\n", name, image)
	}
	if m.cli == "" {
		return "", errors.New("no container CLI found: please install podman or docker")
	}

	// 1. Pull image
	if err := m.runCmdCapture(m.cli, "pull", image); err != nil {
		return "", fmt.Errorf("failed to pull image %s: %w", image, err)
	}

	// 2. Create container and capture returned container ID
	out, err := m.runCmdOutput(m.cli, "create", "--name", name, image, "tail", "-f", "/dev/null")
	if err != nil {
		return "", fmt.Errorf("failed to create container %s from image %s: %w", name, image, err)
	}
	containerID := strings.TrimSpace(out)
	if containerID == "" {
		// Some CLIs may not print the ID; fall back to using the provided name
		containerID = name
	}

	// 3. Start container
	if err := m.runCmdCapture(m.cli, "start", containerID); err != nil {
		// Attempt cleanup: remove container if start failed
		_ = m.runCmdCapture(m.cli, "rm", "-f", containerID)
		return "", fmt.Errorf("failed to start container %s: %w", containerID, err)
	}

	if m.verbose {
		fmt.Printf("Container %s started (via %s)\n", containerID, m.cli)
	}

	return containerID, nil
}

// runCmdCapture runs a command and returns error with stderr/stdout combined on failure.
func (m *Manager) runCmdCapture(name string, args ...string) error {
	if m.verbose {
		fmt.Printf("exec: %s %s\n", name, strings.Join(args, " "))
	}
	cmd := execCommand(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(out.String() + "\n" + stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("command failed: %s %v: %s", name, args, msg)
	}
	return nil
}

// runCmdOutput runs a command and returns its stdout (trimmed). On error, stderr/stdout are included in the error.
func (m *Manager) runCmdOutput(name string, args ...string) (string, error) {
	if m.verbose {
		fmt.Printf("exec: %s %s\n", name, strings.Join(args, " "))
	}
	cmd := execCommand(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(out.String() + "\n" + stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("command failed: %s %v: %s", name, args, msg)
	}
	return out.String(), nil
}

// RunCommand executes a command in a container
func (m *Manager) RunCommand(containerID string, command string) error {
	if m.verbose {
		fmt.Printf("Running command in %s: %s\n", containerID, command)
	}
	if m.cli == "" {
		return errors.New("no container CLI found: please install podman or docker")
	}

	// Use `exec` to run the command inside the container. Use sh -lc to support complex commands.
	// Stream stdout/stderr to the current process so callers see realtime output.
	args := []string{"exec", "-i", containerID, "sh", "-lc", command}
	if m.verbose {
		fmt.Printf("exec: %s %s\n", m.cli, strings.Join(args, " "))
	}

	cmd := execCommand(m.cli, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// No stdin wiring for now; could be added if needed

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to exec command in container %s: %w", containerID, err)
	}

	return nil
}

// RemoveContainer removes a Podman container
func (m *Manager) RemoveContainer(containerID string) error {
	if m.verbose {
		fmt.Printf("Removing container: %s\n", containerID)
	}
	if m.cli == "" {
		return errors.New("no container CLI found: please install podman or docker")
	}

	// Attempt to stop the container; ignore error if it is already stopped
	_ = m.runCmdCapture(m.cli, "stop", containerID)

	// Remove the container forcefully
	if err := m.runCmdCapture(m.cli, "rm", "-f", containerID); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}

	if m.verbose {
		fmt.Printf("Container %s removed\n", containerID)
	}

	return nil
}
