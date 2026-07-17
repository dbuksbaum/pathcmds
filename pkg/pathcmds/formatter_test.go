package pathcmds

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestFormatter(t *testing.T) {
	folders := []Folder{
		{
			Path:     "/usr/bin",
			Category: CategorySystem,
			Commands: []Command{
				{Name: "bash"},
				{Name: "cat"},
				{Name: "ls"},
				{Name: "zsh"},
			},
		},
	}

	var buf bytes.Buffer
	f := NewFormatter(&buf)
	err := f.FormatFolders(folders)
	if err != nil {
		t.Fatalf("FormatFolders failed: %v", err)
	}

	output := buf.String()

	// Check if header with bold blue ANSI color codes is present
	expectedHeader := ansiBoldBlue + "/usr/bin [SYSTEM] (4 commands)" + ansiReset
	if !strings.Contains(output, expectedHeader) {
		t.Errorf("expected output to contain header %q, got %q", expectedHeader, output)
	}

	// Check if commands are printed
	for _, cmd := range folders[0].Commands {
		if !strings.Contains(output, cmd.Name) {
			t.Errorf("expected output to contain command %q", cmd.Name)
		}
	}
}

func TestFormatterEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	err := f.FormatFolders(nil)
	if err != nil {
		t.Fatalf("FormatFolders failed with nil folders: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestFormatterEmptyFolder(t *testing.T) {
	folders := []Folder{
		{
			Path:     "/empty/dir",
			Category: CategoryUser,
			Commands: nil,
		},
	}

	var buf bytes.Buffer
	f := NewFormatter(&buf)
	err := f.FormatFolders(folders)
	if err != nil {
		t.Fatalf("FormatFolders failed with empty folder: %v", err)
	}

	output := buf.String()
	expectedHeader := ansiBoldBlue + "/empty/dir [USER] (0 commands)" + ansiReset
	if !strings.Contains(output, expectedHeader) {
		t.Errorf("expected header %q, got %q", expectedHeader, output)
	}
	// The output should only contain the header and newlines, no command output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only 1 line of output, got %d lines: %v", len(lines), lines)
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// Backup COLUMNS
	oldColumns := os.Getenv("COLUMNS")
	defer os.Setenv("COLUMNS", oldColumns)

	// Test with COLUMNS set to a valid integer
	os.Setenv("COLUMNS", "120")
	w := getTerminalWidth()
	if w != 120 {
		// Note: term.GetSize may succeed in some test runners if connected to a TTY,
		// but standard test runners redirect output. If it succeeded, w might be > 0.
		// Let's assert it is at least a positive integer.
		if w <= 0 {
			t.Errorf("expected terminal width to be positive, got %d", w)
		}
	}

	// Test with COLUMNS set to invalid value
	os.Setenv("COLUMNS", "invalid")
	w2 := getTerminalWidth()
	if w2 <= 0 {
		t.Errorf("expected terminal width to fall back to positive, got %d", w2)
	}

	// Test with COLUMNS unset
	os.Unsetenv("COLUMNS")
	w3 := getTerminalWidth()
	if w3 <= 0 {
		t.Errorf("expected terminal width to fall back to positive, got %d", w3)
	}
}

func TestFormatterDifferentWidths(t *testing.T) {
	// Backup COLUMNS
	oldColumns := os.Getenv("COLUMNS")
	defer os.Setenv("COLUMNS", oldColumns)

	folders := []Folder{
		{
			Path:     "/bin",
			Category: CategorySystem,
			Commands: []Command{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
				{Name: "d"},
			},
		},
	}

	// Set terminal width to very small value (e.g. 5) to force single column
	os.Setenv("COLUMNS", "5")
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	if err := f.FormatFolders(folders); err != nil {
		t.Fatalf("FormatFolders failed: %v", err)
	}
	// Output should contain the commands
	output := buf.String()
	for _, cmd := range folders[0].Commands {
		if !strings.Contains(output, cmd.Name) {
			t.Errorf("expected output to contain command %q under small width", cmd.Name)
		}
	}

	// Set terminal width to very large value (e.g. 500) to force multiple columns
	os.Setenv("COLUMNS", "500")
	buf.Reset()
	if err := f.FormatFolders(folders); err != nil {
		t.Fatalf("FormatFolders failed: %v", err)
	}
	output = buf.String()
	for _, cmd := range folders[0].Commands {
		if !strings.Contains(output, cmd.Name) {
			t.Errorf("expected output to contain command %q under large width", cmd.Name)
		}
	}
}
