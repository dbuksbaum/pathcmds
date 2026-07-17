package pathcmds

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var (
	// execLookPath is a variable wrapper around exec.LookPath to facilitate unit testing.
	execLookPath = exec.LookPath
	// execCommand is a variable wrapper around exec.Command to facilitate unit testing.
	execCommand  = exec.Command
)

// Pager manages spawning the 'less -R' paging sub-process and piping output into it.
type Pager struct {
	cmd    *exec.Cmd      // cmd represents the spawned external pager process ('less')
	stdin  io.WriteCloser // stdin is the write pipe connected to the pager's standard input
	active bool           // active is true if the external pager was successfully started
}

// NewPager configures the output destination. If active is true, it attempts to set up
// the 'less -R' process. If 'less' is missing or active is false, it returns os.Stdout.
func NewPager(active bool) (*Pager, io.Writer, error) {
	if !active {
		return &Pager{active: false}, os.Stdout, nil
	}

	// Verify 'less' binary is available in the environment's executable path
	lessPath, err := execLookPath("less")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: 'less' binary not found. Falling back to stdout without paging.")
		return &Pager{active: false}, os.Stdout, nil
	}

	// Spawn 'less -R' to handle ANSI escape sequence colors correctly
	cmd := execCommand(lessPath, "-R")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stdin pipe for less: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start less sub-process: %w", err)
	}

	return &Pager{
		cmd:    cmd,
		stdin:  stdin,
		active: true,
	}, stdin, nil
}

// Close closes the pager's pipe and waits for the user to exit the less interface.
func (p *Pager) Close() error {
	if !p.active {
		return nil
	}

	// Close the stdin pipe to notify 'less' that writing has finished
	_ = p.stdin.Close()

	// Wait for the less pager to finish (i.e. until the user presses 'q' to quit)
	_ = p.cmd.Wait()

	return nil
}
