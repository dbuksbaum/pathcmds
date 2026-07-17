package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// captureStdout executes a function and returns whatever was written to os.Stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	f()

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func resetFlags() {
	flagSystem = false
	flagUser = false
	flagApps = false
	flagPage = false

	RootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})
}

func TestVersionCommand(t *testing.T) {
	resetFlags()
	RootCmd.SetArgs([]string{"version"})

	output := captureStdout(t, func() {
		err := RootCmd.Execute()
		if err != nil {
			t.Errorf("version execution failed: %v", err)
		}
	})

	if !strings.Contains(output, "pathcmds v1.0.0") {
		t.Errorf("expected version output 'pathcmds v1.0.0', got %q", output)
	}
}

func TestInvalidCommand(t *testing.T) {
	// Set PATH to a non-existent path
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", "/nonexistent/test/path/for/invalid/cmd")

	resetFlags()
	RootCmd.SetArgs([]string{"invalid"})

	output := captureStdout(t, func() {
		err := RootCmd.Execute()
		if err != nil {
			t.Errorf("invalid command execution failed: %v", err)
		}
	})

	if !strings.Contains(output, "Found 1 invalid PATH entries:") {
		t.Errorf("expected count warning, got %q", output)
	}
	if !strings.Contains(output, "does not exist") {
		t.Errorf("expected 'does not exist' reason, got %q", output)
	}
}

func TestRootCommandHelp(t *testing.T) {
	resetFlags()
	RootCmd.SetArgs([]string{"--help"})

	output := captureStdout(t, func() {
		_ = RootCmd.Execute()
	})

	if !strings.Contains(output, "pathcmds is a CLI utility that parses your $PATH") {
		t.Errorf("expected help output, got %q", output)
	}
}

func TestRootCommandExecution(t *testing.T) {
	// Setup a temp directory with a mock command
	tempDir, err := os.MkdirTemp("", "pathcmds-cmd-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an executable file inside tempDir (e.g. mock-cli)
	execPath := filepath.Join(tempDir, "mock-cli")
	err = os.WriteFile(execPath, []byte("#!/bin/sh\necho mock"), 0755)
	if err != nil {
		t.Fatalf("failed to create exec file: %v", err)
	}

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tempDir)

	// Test default run (no flags)
	resetFlags()
	RootCmd.SetArgs([]string{})

	output := captureStdout(t, func() {
		err := RootCmd.Execute()
		if err != nil {
			t.Errorf("root cmd execution failed: %v", err)
		}
	})

	if !strings.Contains(output, "mock-cli") {
		t.Errorf("expected output to contain command 'mock-cli', got %q", output)
	}
	// By default, it categorizes custom paths as CategoryUser
	if !strings.Contains(output, "[USER]") {
		t.Errorf("expected category header '[USER]', got %q", output)
	}

	// Test with system flag (should not match mock-cli since it categorizes as User by default)
	resetFlags()
	RootCmd.SetArgs([]string{"--system"})

	output = captureStdout(t, func() {
		err := RootCmd.Execute()
		if err != nil {
			t.Errorf("root cmd execution with --system failed: %v", err)
		}
	})

	if strings.Contains(output, "mock-cli") {
		t.Errorf("expected output to NOT contain 'mock-cli' when filtered by --system, got %q", output)
	}

	// Test with user flag (should match mock-cli)
	resetFlags()
	RootCmd.SetArgs([]string{"-u"})

	output = captureStdout(t, func() {
		err := RootCmd.Execute()
		if err != nil {
			t.Errorf("root cmd execution with -u failed: %v", err)
		}
	})

	if !strings.Contains(output, "mock-cli") {
		t.Errorf("expected output to contain 'mock-cli' when filtered by -u, got %q", output)
	}
}
