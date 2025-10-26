package container

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// fakeExec simulates the container CLI for tests.
// It inspects the combined command (name + args) and returns a shell command that
// prints appropriate outputs for create/exec and otherwise exits zero.
func fakeExec(name string, args ...string) *exec.Cmd {
	full := strings.Join(append([]string{name}, args...), " ")
	var script string
	switch {
	case strings.Contains(full, " create "):
		// simulate returning a container id
		script = "printf 'fake-id-123'"
	case strings.Contains(full, " exec "):
		// simulate command output
		script = "echo hello"
	default:
		// simulate success for pull/start/stop/rm
		script = "exit 0"
	}
	return exec.Command("sh", "-c", script)
}

func TestCreateContainer_ReturnsID(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExec

	m := NewManager(false)
	m.cli = "podman"

	id, err := m.CreateContainer("ubuntu:22.04", "testname")
	if err != nil {
		t.Fatalf("CreateContainer failed: %v", err)
	}
	if id != "fake-id-123" {
		t.Fatalf("expected fake-id-123, got %q", id)
	}
}

func TestRunCommand_StreamsOutput(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExec

	m := NewManager(false)
	m.cli = "podman"

	// capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := m.RunCommand("fake-id", "echo hello")
	// restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := strings.TrimSpace(buf.String())
	if out != "hello" {
		t.Fatalf("expected output 'hello', got %q", out)
	}
}

func TestRemoveContainer_Succeeds(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExec

	m := NewManager(false)
	m.cli = "podman"

	if err := m.RemoveContainer("fake-id"); err != nil {
		t.Fatalf("RemoveContainer failed: %v", err)
	}
}

func TestPullImage_Succeeds(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExec

	m := NewManager(false)
	m.cli = "podman"

	if err := m.PullImage("ubuntu:22.04"); err != nil {
		t.Fatalf("PullImage failed: %v", err)
	}
}

// fakeFailExec simulates a failing pull command.
func fakeFailExec(name string, args ...string) *exec.Cmd {
	full := strings.Join(append([]string{name}, args...), " ")
	if strings.Contains(full, " pull ") {
		// return a command that prints an error message and exits non-zero
		return exec.Command("sh", "-c", "echo failed pull >&2; exit 2")
	}
	// fallback to success for other commands
	return exec.Command("sh", "-c", "exit 0")
}

func TestPullImage_Failure(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeFailExec

	m := NewManager(false)
	m.cli = "podman"

	if err := m.PullImage("nonexistent:image"); err == nil {
		t.Fatalf("expected PullImage to fail for nonexistent image")
	}
}
