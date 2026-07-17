package pathcmds

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "pathcmds-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an executable file
	execPath := filepath.Join(tempDir, "b-exec")
	err = os.WriteFile(execPath, []byte("#!/bin/sh\necho b"), 0755) // executable
	if err != nil {
		t.Fatalf("failed to create exec file: %v", err)
	}

	// Create another executable file to test case-insensitive sorting
	execPath2 := filepath.Join(tempDir, "A-exec")
	err = os.WriteFile(execPath2, []byte("#!/bin/sh\necho a"), 0755) // executable
	if err != nil {
		t.Fatalf("failed to create exec file: %v", err)
	}

	// Create a non-executable file
	nonExecPath := filepath.Join(tempDir, "my-non-exec")
	err = os.WriteFile(nonExecPath, []byte("some text"), 0644) // non-executable
	if err != nil {
		t.Fatalf("failed to create non-exec file: %v", err)
	}

	// Create a subdirectory inside tempDir (should be skipped)
	subDirPath := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDirPath, 0755)
	if err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// Set PATH to tempDir
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tempDir)

	parser := NewParser()
	folders, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(folders) != 1 {
		t.Fatalf("expected 1 folder, got %d", len(folders))
	}

	folder := folders[0]
	// Clean/resolve the temp directory for absolute path matching
	expectedPath, err := filepath.Abs(tempDir)
	if err != nil {
		expectedPath = tempDir
	}

	if folder.Path != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, folder.Path)
	}

	// Expected commands count should be 2: A-exec and b-exec
	if len(folder.Commands) != 2 {
		t.Fatalf("expected 2 commands in folder, got %d", len(folder.Commands))
	}

	// Verify case-insensitive sorting (A-exec should come before b-exec)
	if folder.Commands[0].Name != "A-exec" {
		t.Errorf("expected first command to be 'A-exec' (sorted), got %q", folder.Commands[0].Name)
	}
	if folder.Commands[1].Name != "b-exec" {
		t.Errorf("expected second command to be 'b-exec', got %q", folder.Commands[1].Name)
	}
}

func TestParseEmptyPATH(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", "")

	parser := NewParser()
	_, err := parser.Parse()
	if err == nil {
		t.Error("expected error for empty PATH, got nil")
	}
}

func TestParseDetailedInvalidPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pathcmds-invalid-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file to use as a path entry (not a directory)
	filePath := filepath.Join(tempDir, "regular-file")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// Create a directory with no permissions (locked)
	lockedDir := filepath.Join(tempDir, "locked-dir")
	if err := os.Mkdir(lockedDir, 0000); err != nil {
		t.Fatalf("failed to create locked dir: %v", err)
	}
	// Restore permissions on defer so RemoveAll can delete it
	defer os.Chmod(lockedDir, 0755)

	nonExistentDir := filepath.Join(tempDir, "does-not-exist")

	// Set PATH to contain: nonExistentDir, filePath, lockedDir
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", strings.Join([]string{nonExistentDir, filePath, lockedDir}, string(filepath.ListSeparator)))

	parser := NewParser()
	folders, invalidPaths, err := parser.ParseDetailed()
	if err != nil {
		t.Fatalf("ParseDetailed failed: %v", err)
	}

	if len(folders) != 0 {
		t.Errorf("expected 0 folders, got %d", len(folders))
	}

	// We expect 3 invalid paths (nonExistentDir, filePath, lockedDir)
	if len(invalidPaths) != 3 {
		t.Fatalf("expected 3 invalid paths, got %d", len(invalidPaths))
	}

	reasons := make(map[string]string)
	for _, ip := range invalidPaths {
		reasons[filepath.Base(ip.Path)] = ip.Reason
	}

	if r, ok := reasons["does-not-exist"]; !ok || !strings.Contains(r, "does not exist") {
		t.Errorf("expected reason for 'does-not-exist' to contain 'does not exist', got %q", r)
	}
	if r, ok := reasons["regular-file"]; !ok || !strings.Contains(r, "not a directory") {
		t.Errorf("expected reason for 'regular-file' to be 'not a directory', got %q", r)
	}
	if r, ok := reasons["locked-dir"]; !ok || (!strings.Contains(r, "permission denied") && !strings.Contains(r, "locked")) {
		t.Errorf("expected reason for 'locked-dir' to contain 'permission denied' or 'locked', got %q", r)
	}
}

func TestParseWarningsToStderr(t *testing.T) {
	// Temporarily capture os.Stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	// Set PATH to a non-existent directory
	oldPath := os.Getenv("PATH")
	defer func() {
		os.Setenv("PATH", oldPath)
		os.Stderr = oldStderr
	}()
	os.Setenv("PATH", "/nonexistent/path/for/warnings/test")

	parser := NewParser()
	_, err = parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Close writer so we can read from reader without blocking
	w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Warning: path") || !strings.Contains(output, "does not exist") {
		t.Errorf("expected warning printed to stderr, got output: %q", output)
	}
}
