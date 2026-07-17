package pathcmds

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Categories defined for CLI filters
const (
	CategorySystem = "system"
	CategoryUser   = "user"
	CategoryApp    = "app"
)

// Categorizer implements the classification rules for PATH directories.
type Categorizer struct {
	appRegex *regexp.Regexp
}

// NewCategorizer instantiates a Categorizer with predefined rules.
func NewCategorizer() *Categorizer {
	// App pattern matches specific toolchains, SDKs, applications, or runtimes.
	// Examples: .cargo/bin, .rustup, .dotnet/sdk, nvm, node, go/bin, .rvm, .pyenv, Applications, etc.
	appPattern := `(?i)(\.dotnet|\.cargo|\.rustup|/go/bin|\.nvm|\.n/bin|/node/|pnpm|yarn|\.rvm|\.pyenv|/Library/Android/|/Applications/|/Library/Developer/|/sdk/|/toolchains/|\.local/share/|cellar)`
	
	return &Categorizer{
		appRegex: regexp.MustCompile(appPattern),
	}
}

// Categorize determines whether a folder path falls into "system", "user", or "app".
func (c *Categorizer) Categorize(path string) string {
	cleaned := filepath.Clean(path)
	lower := strings.ToLower(cleaned)

	// Homebrew paths are explicitly user command sets
	if strings.Contains(lower, "homebrew") {
		return CategoryUser
	}

	// 1. App category takes priority (e.g. if inside a user home folder like /Users/username/.cargo/bin, or vendor opt paths)
	if c.appRegex.MatchString(cleaned) || strings.HasPrefix(lower, "/opt/") {
		return CategoryApp
	}

	// 2. System category matching
	// Pure system commands should reside in /bin, /sbin, /usr/bin, /usr/sbin.
	// Exclude user-writable system spaces like /usr/local/bin.
	if lower == "/bin" || lower == "/sbin" ||
		lower == "/usr/bin" || lower == "/usr/sbin" ||
		strings.HasPrefix(lower, "/bin/") || strings.HasPrefix(lower, "/sbin/") ||
		strings.HasPrefix(lower, "/usr/bin/") || strings.HasPrefix(lower, "/usr/sbin/") ||
		strings.HasPrefix(lower, "/system/") {
		return CategorySystem
	}

	// 3. User category matching
	// Includes /usr/local/bin, /usr/local/sbin, homebrew (/opt/homebrew), and paths in user home directories.
	if strings.Contains(lower, "/usr/local") ||
		strings.Contains(lower, "/opt/homebrew") ||
		strings.Contains(lower, "/home/") ||
		strings.Contains(lower, "/users/") ||
		strings.HasPrefix(lower, "~") ||
		strings.HasSuffix(lower, "/bin") ||
		strings.HasSuffix(lower, "/sbin") ||
		strings.HasSuffix(lower, "/.local/bin") {
		return CategoryUser
	}

	// Default fallback to user category for custom paths
	return CategoryUser
}
