package pathcmds

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"golang.org/x/term"
)

const (
	ansiBoldBlue = "\033[1;34m"
	ansiReset    = "\033[0m"
)

// Formatter handles the grid-aligned layout and coloring of directories and their commands.
type Formatter struct {
	writer io.Writer
}

// NewFormatter returns a new Formatter that writes to the given io.Writer.
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{writer: w}
}

// FormatFolders outputs the given folders and their commands as a colored list with a column-aligned grid.
func (f *Formatter) FormatFolders(folders []Folder) error {
	termWidth := getTerminalWidth()

	for i, folder := range folders {
		if i > 0 {
			// Insert a separating blank line between directory groups
			fmt.Fprintln(f.writer)
		}

		// Bold, blue header format: PATH [CATEGORY] (COUNT commands)
		header := fmt.Sprintf("%s%s [%s] (%d commands)%s",
			ansiBoldBlue,
			folder.Path,
			strings.ToUpper(folder.Category),
			len(folder.Commands),
			ansiReset,
		)
		fmt.Fprintln(f.writer, header)

		if len(folder.Commands) == 0 {
			continue
		}

		// Determine maximum command name length to calculate column sizing
		maxLen := 0
		for _, cmd := range folder.Commands {
			if len(cmd.Name) > maxLen {
				maxLen = len(cmd.Name)
			}
		}

		// Determine grid columns
		padding := 4
		colWidth := maxLen + padding
		numCols := termWidth / colWidth
		if numCols <= 0 {
			numCols = 1
		}

		numRows := (len(folder.Commands) + numCols - 1) / numCols

		// Setup tabwriter with space-padding and standard settings
		tw := tabwriter.NewWriter(f.writer, 0, 0, padding, ' ', 0)

		// Column-major layout loop
		for r := 0; r < numRows; r++ {
			var lineParts []string
			for c := 0; c < numCols; c++ {
				idx := c*numRows + r
				if idx < len(folder.Commands) {
					lineParts = append(lineParts, folder.Commands[idx].Name)
				}
			}
			// Each column needs to end with a tab character to notify tabwriter of boundary
			fmt.Fprintln(tw, strings.Join(lineParts, "\t")+"\t")
		}

		if err := tw.Flush(); err != nil {
			return fmt.Errorf("failed to flush tabwriter: %w", err)
		}
	}

	return nil
}

// getTerminalWidth determines the terminal width by querying stdout/stderr/stdin descriptors or the COLUMNS env var.
func getTerminalWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	if w, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil && w > 0 {
		return w
	}
	if w, _, err := term.GetSize(int(os.Stdin.Fd())); err == nil && w > 0 {
		return w
	}
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var w int
		if _, err := fmt.Sscan(cols, &w); err == nil && w > 0 {
			return w
		}
	}
	return 80 // Default to 80 if no terminal width is detected
}
