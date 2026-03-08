package container

import (
	"os/exec"
	"strings"
	"testing"
)

// fakeExecConfig simulates the container CLI for config-aware tests.
// It returns a container id for create and succeeds for other commands.
func fakeExecConfig(name string, args ...string) *exec.Cmd {
	full := strings.Join(append([]string{name}, args...), " ")
	switch {
	case strings.Contains(full, " create "):
		// simulate returning a container id
		return exec.Command("sh", "-c", "printf 'fake-id-123'")
	default:
		return exec.Command("sh", "-c", "exit 0")
	}
}

func TestCreateContainerWithConfig_ReturnsID(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExecConfig

	m := NewManager(false)
	m.cli = "podman"

	cfg := &ContainerConfig{
		Env:     []string{"FOO=bar", "BAZ=quux"},
		Volumes: []string{"/host:/container:ro"},
		WorkDir: "/container",
		User:    "1000:1000",
	}

	id, err := m.CreateContainerWithConfig("ubuntu:22.04", "testname", cfg)
	if err != nil {
		t.Fatalf("CreateContainerWithConfig failed: %v", err)
	}
	if id != "fake-id-123" {
		t.Fatalf("expected fake-id-123, got %q", id)
	}
}
