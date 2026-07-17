# Tutorial Guide: Getting Started with pathcmds

This tutorial will guide you step-by-step through installing, running, and using `pathcmds` to inspect, clean, and manage your command-line environment's `$PATH` variable.

---

## Prerequisites

Before beginning, ensure you have:
- **Go** (version 1.26.5 or later) installed on your system.
- A terminal (Bash, Zsh, or Fish) on macOS, Linux, or another Unix-like system.
- `just` (optional, but recommended for task execution).

---

## Step 1: Clone and Build the Project

First, let's clone the repository and compile the executable.

```bash
# Clone the repository
git clone https://github.com/davidbuksbaum/pathcmds.git
cd pathcmds

# Build the executable using just
just build
```

If you don't have `just` installed, you can initialize modules and build using standard Go commands:
```bash
go mod tidy
mkdir -p bin
go build -o bin/pathcmds main.go
```

This creates a self-contained executable at `bin/pathcmds`.

---

## Step 2: Run Your First Inspection

Let's see what is currently on your `$PATH`. Run `pathcmds` without any flags:

```bash
./bin/pathcmds
```

You will see output formatted like this:
```ansi
/usr/bin [SYSTEM] (45 commands)
bash      cat       cp        dd        df        du        find
grep      gzip      kill      ln        ls        mkdir     mv
ps        pwd       rm        rmdir     sed       tar       zsh
...

/usr/local/bin [USER] (12 commands)
docker             git                node               npm
...
```

Notice:
1. Directory paths are in **bold blue** (if your terminal supports ANSI color).
2. The category (`[SYSTEM]`, `[USER]`, or `[APP]`) is displayed next to each folder.
3. The number of commands found is listed.
4. Commands are laid out in a column-aligned grid matching your current terminal width.

---

## Step 3: Filter by Categories

In most development workflows, you only care about commands you installed recently (user/application tools) rather than default system commands. You can filter the display using categories.

### View Only System Commands
To inspect standard system commands located in `/bin`, `/sbin`, `/usr/bin`, and `/usr/sbin`:
```bash
./bin/pathcmds --system
```

### View User-installed Tools and Custom Binaries
To see tools installed via Homebrew, Node packages, or manual user folder bins:
```bash
./bin/pathcmds --user
```

### View Application SDKs and Runtimes
To see executables inside package managers and runtimes (like Rust Cargo, Go, .NET, Python pyenv, Node nvm, etc.):
```bash
./bin/pathcmds --apps
```

### Combine Filters
You can combine flags to see multiple categories together. For example, to view both your User and App categories:
```bash
./bin/pathcmds -u -a
```

---

## Step 4: Page Large Output

If you have a very long `$PATH` variable, listing all commands can clutter your terminal history. `pathcmds` has a built-in paging mechanism that feeds output into the system's `less` pager while maintaining color styling:

```bash
./bin/pathcmds --page
# Or use the shorthand
./bin/pathcmds -p
```

Inside the pager:
- Use `Up/Down` arrow keys or `PageUp/PageDown` to scroll.
- Press `/` followed by a keyword to search for a command.
- Press `q` to exit the pager and return to your terminal prompt.

---

## Step 5: Find and Fix Broken PATH Entries

Over time, installing and uninstalling software can leave dead, non-existent, or locked directory references in your `$PATH`. This slows down shell startup and causes shell autocomplete lag.

Let's check if you have any invalid entries:

```bash
./bin/pathcmds invalid
```

### Case A: Healthy PATH
If your configuration is clean, you will see:
```
All entries in $PATH are valid.
```

### Case B: Broken PATH
If there are errors, they will be reported:
```
Found 2 invalid PATH entries:
- /usr/local/opt/oldtool/bin: does not exist
- /Users/john/.config/private: permission denied / locked (stat: permission denied)
```

### How to Clean Up Broken Paths:
1. Open your shell configuration file in a text editor (e.g. `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`).
2. Search for the reported paths (e.g. `oldtool` or `private`).
3. Remove or comment out the lines exporting those paths, or fix the typo.
4. Restart your shell (or run `source ~/.zshrc`), then run `./bin/pathcmds invalid` again to verify!

---

## Summary

You are now ready to use `pathcmds` as a regular tool in your terminal toolkit. You can install it globally to make it accessible everywhere:

```bash
just install
```

Happy command-line tuning!
