// Package pathcmds provides functionality to parse, categorize, format, and page paths.
package pathcmds

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Command represents an executable command found in a directory.
type Command struct {
	Name string
}

// Folder represents a directory on the user's PATH containing its scanned executables.
type Folder struct {
	Path     string
	Category string
	Commands []Command
}

// Parser manages the extraction and validation of directories in the PATH.
type Parser struct{}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// InvalidPath represents an entry on the user's PATH that is broken, missing, or inaccessible.
type InvalidPath struct {
	Path   string
	Reason string
}

// ParseDetailed retrieves the user's PATH environment variable, splits it,
// and returns successfully parsed folders along with details of invalid paths.
func (p *Parser) ParseDetailed() ([]Folder, []InvalidPath, error) {
	pathVal := os.Getenv("PATH")
	if pathVal == "" {
		return nil, nil, fmt.Errorf("PATH environment variable is empty or not set")
	}

	paths := filepath.SplitList(pathVal)
	var folders []Folder
	var invalidPaths []InvalidPath

	for _, path := range paths {
		if path == "" {
			continue
		}

		// Resolve clean absolute path for consistency
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = path // fallback to original path if Abs fails
		}

		info, err := os.Stat(absPath)
		if err != nil {
			reason := "does not exist"
			if !os.IsNotExist(err) {
				reason = fmt.Sprintf("inaccessible (%v)", err)
			}
			invalidPaths = append(invalidPaths, InvalidPath{Path: absPath, Reason: reason})
			continue
		}

		if !info.IsDir() {
			invalidPaths = append(invalidPaths, InvalidPath{Path: absPath, Reason: "not a directory"})
			continue
		}

		// Read the directory contents
		entries, err := os.ReadDir(absPath)
		if err != nil {
			invalidPaths = append(invalidPaths, InvalidPath{Path: absPath, Reason: fmt.Sprintf("permission denied / locked (%v)", err)})
			continue
		}

		var commands []Command
		for _, entry := range entries {
			// Skip directories inside PATH folders
			if entry.IsDir() {
				continue
			}

			// Retrieve file info to check executable bit
			fileInfo, err := entry.Info()
			if err != nil {
				continue
			}

			if isExecutable(fileInfo) {
				commands = append(commands, Command{Name: entry.Name()})
			}
		}

		// Filter out empty folders
		if len(commands) == 0 {
			continue
		}

		// Sort commands alphabetically (case-insensitive)
		sort.Slice(commands, func(i, j int) bool {
			return strings.ToLower(commands[i].Name) < strings.ToLower(commands[j].Name)
		})

		folders = append(folders, Folder{
			Path:     absPath,
			Commands: commands,
		})
	}

	return folders, invalidPaths, nil
}

// Parse retrieves the user's PATH environment variable, splits it, and scans each directory.
// It skips missing directories, files that are directories, and non-executable files.
// Broken or inaccessible paths will print a warning to stderr instead of causing a crash.
func (p *Parser) Parse() ([]Folder, error) {
	folders, invalidPaths, err := p.ParseDetailed()
	if err != nil {
		return nil, err
	}
	for _, ip := range invalidPaths {
		fmt.Fprintf(os.Stderr, "Warning: path %q %s (skipping)\n", ip.Path, ip.Reason)
	}
	return folders, nil
}

// isExecutable checks if the file has executable permissions set on Unix systems.
func isExecutable(info fs.FileInfo) bool {
	// Mode bit 0111 checks for user, group, and other execute permission
	return info.Mode()&0111 != 0
}
