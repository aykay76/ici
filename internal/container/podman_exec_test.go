package container

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// fakeExecStreams prints provided output for exec commands
func fakeExecStreams(name string, args ...string) *exec.Cmd {
	full := strings.Join(append([]string{name}, args...), " ")
	if strings.Contains(full, " exec ") {
		// print both stdout and stderr and exit 0
		script := "echo out; echo err 1>&2"
		return exec.Command("sh", "-c", script)
	}
	// default no-op
	return exec.Command("sh", "-c", "exit 0")
}

func TestRunCommandWithOptions_StreamsOutput(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExecStreams

	m := NewManager(false)
	m.cli = "podman"

	// capture stdout
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	err := m.RunCommandWithOptions("fake-id", "echo hello", nil, "", 0)

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("RunCommandWithOptions failed: %v", err)
	}

	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	_, _ = io.Copy(&bufOut, rOut)
	_, _ = io.Copy(&bufErr, rErr)

	if strings.TrimSpace(bufOut.String()) != "out" {
		t.Fatalf("expected stdout 'out', got %q", bufOut.String())
	}
	if strings.TrimSpace(bufErr.String()) != "err" {
		t.Fatalf("expected stderr 'err', got %q", bufErr.String())
	}
}

// fakeExecInspect returns the command string so tests can assert flags were added
func fakeExecInspect(name string, args ...string) *exec.Cmd {
	// print args to stdout so we can inspect them
	script := "printf '%s' \"" + strings.Join(append([]string{name}, args...), " ") + "\""
	return exec.Command("sh", "-c", script)
}

func TestRunCommandWithOptions_EnvAndWorkdir(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExecInspect

	m := NewManager(false)
	m.cli = "podman"

	// capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := m.RunCommandWithOptions("fake-id", "echo hello", []string{"FOO=bar"}, "/workspace/app", 0)

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunCommandWithOptions failed: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "--env FOO=bar") {
		t.Fatalf("expected --env to be present in exec args, got %q", out)
	}
	if !strings.Contains(out, "-w /workspace/app") {
		t.Fatalf("expected -w workdir to be present in exec args, got %q", out)
	}
}

// fakeExecSleep sleeps for longer than the timeout
func fakeExecSleep(name string, args ...string) *exec.Cmd {
	if strings.Contains(strings.Join(args, " "), "exec") {
		return exec.Command("sh", "-c", "sleep 5; exit 0")
	}
	return exec.Command("sh", "-c", "exit 0")
}

func TestRunCommandWithOptions_Timeout(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()
	execCommand = fakeExecSleep

	m := NewManager(false)
	m.cli = "podman"

	start := time.Now()
	err := m.RunCommandWithOptions("fake-id", "sleep 5", nil, "", 1)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
	if elapsed > 3*time.Second {
		t.Fatalf("timeout did not occur promptly, elapsed: %v", elapsed)
	}
}
