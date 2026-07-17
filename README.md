# pathcmds

[![Go Reference](https://pkg.go.dev/badge/pathcmds.svg)](https://pkg.go.dev/pathcmds)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidbuksbaum/pathcmds)](https://goreportcard.com/report/github.com/davidbuksbaum/pathcmds)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`pathcmds` is a fast, lightweight, and modern CLI utility written in Go to inspect, categorize, and format your shell's `$PATH` environment variable. It helps you quickly discover what executables are available, group them by folder, and identify invalid, duplicate, or broken paths.

---

## Features

- **Categorized View**: Groups `$PATH` directories into `SYSTEM`, `USER`, or `APP` sets for easier visualization.
- **Smart Formatting**: Displays commands in a terminal-width-aware, column-aligned grid.
- **Colorized Headers**: Distinguishes directory headers with clean ANSI colors.
- **Built-in Pagination**: Seamlessly pages long command lists using your system's `less -R` when requested.
- **Diagnostics**: Instantly identifies broken, non-existent, or permission-restricted directories with the `invalid` command.

---

## Installation

### From Source

Ensure you have Go (1.26.5 or later) installed.

Clone the repository and build using the provided `justfile` or manually:

```bash
# Clone the repository
git clone https://github.com/davidbuksbaum/pathcmds.git
cd pathcmds

# Build the executable (placed in bin/pathcmds)
just build

# Install the executable to your GOBIN directory
just install
```

---

## Usage

Run `pathcmds` without arguments to inspect and list all executables grouped by directory:

```bash
pathcmds
```

### Filtering Categories

You can filter command sets by passing category flags:

- `-s, --system`: Filter and show system command sets (e.g., `/bin`, `/usr/bin`, `/sbin`, `/usr/sbin`).
- `-u, --user`: Filter and show user-configured command sets (e.g., `/usr/local/bin`, `/opt/homebrew/bin`, `~/.local/bin`).
- `-a, --apps`: Filter and show application-specific folders (e.g., `.dotnet/tools`, `.cargo/bin`, `.nvm/versions/node/...`).

```bash
# Show only system executables
pathcmds --system

# Show both user and app-specific executables
pathcmds -u -a
```

### Paging

For very large `$PATH` configurations, use the `-p` / `--page` flag to pipe the output through the system's `less` pager (retaining ANSI colors):

```bash
pathcmds -p
```

### Diagnostics

Find invalid, missing, or locked entries on your `$PATH`:

```bash
pathcmds invalid
```

If everything is healthy, it returns:
`All entries in $PATH are valid.`

Otherwise, it lists the details of the broken paths:
```
Found 2 invalid PATH entries:
- /usr/local/dummybin: does not exist
- /Users/user/.local/locked: permission denied / locked (stat: permission denied)
```

---

## Development

The project comes with a `justfile` for task management.

- **Initialize dependencies**: `just init`
- **Build the executable**: `just build`
- **Install globally**: `just install`
- **Uninstall**: `just uninstall`
- **Clean build artifacts**: `just clean`
- **Run immediately**: `just run` (e.g. `just run "--system --page"`)

### Running Tests

Execute the unit test suite and verify test coverage:

```bash
go test -v ./...
```

For coverage analysis:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
