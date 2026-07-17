package pathcmds

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
)

// TestNewPagerInactive verifies that when the pager is not active,
// it returns os.Stdout and doesn't spawn any process.
func TestNewPagerInactive(t *testing.T) {
	pager, writer, err := NewPager(false)
	if err != nil {
		t.Fatalf("NewPager(false) failed: %v", err)
	}
	if pager.active {
		t.Error("expected pager to be inactive")
	}
	if writer != os.Stdout {
		t.Error("expected writer to be os.Stdout")
	}

	// Closing inactive pager should succeed without error
	if err := pager.Close(); err != nil {
		t.Errorf("pager.Close() on inactive pager failed: %v", err)
	}
}

// TestNewPagerMissingBinary verifies the fallback behavior to os.Stdout
// when the 'less' binary is not found on the system.
func TestNewPagerMissingBinary(t *testing.T) {
	oldLookPath := execLookPath
	defer func() { execLookPath = oldLookPath }()

	// Mock LookPath to simulate missing 'less' binary
	execLookPath = func(file string) (string, error) {
		return "", errors.New("exec: file not found in $PATH")
	}

	pager, writer, err := NewPager(true)
	if err != nil {
		t.Fatalf("NewPager(true) with missing binary failed: %v", err)
	}
	if pager.active {
		t.Error("expected pager to be inactive when binary is missing")
	}
	if writer != os.Stdout {
		t.Error("expected writer to be os.Stdout when binary is missing")
	}

	if err := pager.Close(); err != nil {
		t.Errorf("pager.Close() failed: %v", err)
	}
}

// TestNewPagerActive verifies that when the pager is active,
// it successfully spawns a process and pipes input to it.
func TestNewPagerActive(t *testing.T) {
	oldLookPath := execLookPath
	oldCommand := execCommand
	defer func() {
		execLookPath = oldLookPath
		execCommand = oldCommand
	}()

	execLookPath = func(file string) (string, error) {
		return "/mock/less", nil
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		// Launch the test binary itself, running the helper process
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	pager, writer, err := NewPager(true)
	if err != nil {
		t.Fatalf("NewPager(true) failed: %v", err)
	}
	if !pager.active {
		t.Fatal("expected pager to be active")
	}

	// Write text to the pager's input pipe
	testStr := "test output for pager\n"
	n, err := io.WriteString(writer, testStr)
	if err != nil {
		t.Fatalf("failed to write to pager: %v", err)
	}
	if n != len(testStr) {
		t.Errorf("wrote %d bytes, expected %d", n, len(testStr))
	}

	// Close the pager and wait for the subprocess to finish
	if err := pager.Close(); err != nil {
		t.Errorf("pager.Close() failed: %v", err)
	}
}

// TestHelperProcess is a helper command used to mock the pager process in tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Copy stdin to stdout and exit cleanly
	_, _ = io.Copy(os.Stdout, os.Stdin)
	os.Exit(0)
}
