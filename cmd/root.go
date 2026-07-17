package cmd

import (
	"fmt"
	"os"

	"pathcmds/pkg/pathcmds"

	"github.com/spf13/cobra"
)

var (
	flagSystem bool
	flagUser   bool
	flagApps   bool
	flagPage   bool
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "pathcmds",
	Short: "Inspect, categorize, and format the user's $PATH environment variable",
	Long: `pathcmds is a CLI utility that parses your $PATH environment variable,
inspects all directories, and retrieves their executable commands.
It groups commands by directory, sorts them alphabetically, and categorizes
each directory into 'system', 'user', or 'app'.

If no category flags (-s, -u, -a) are provided, the tool defaults to displaying
all categories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no category filters are specified, default to enabling all of them
		if !flagSystem && !flagUser && !flagApps {
			flagSystem = true
			flagUser = true
			flagApps = true
		}

		// 1. Parse PATH and read directories
		parser := pathcmds.NewParser()
		folders, err := parser.Parse()
		if err != nil {
			return fmt.Errorf("failed to parse PATH: %w", err)
		}

		// 2. Categorize and filter folders
		categorizer := pathcmds.NewCategorizer()
		var filtered []pathcmds.Folder

		for _, folder := range folders {
			cat := categorizer.Categorize(folder.Path)
			folder.Category = cat

			keep := false
			switch cat {
			case pathcmds.CategorySystem:
				if flagSystem {
					keep = true
				}
			case pathcmds.CategoryUser:
				if flagUser {
					keep = true
				}
			case pathcmds.CategoryApp:
				if flagApps {
					keep = true
				}
			}

			if keep {
				filtered = append(filtered, folder)
			}
		}

		// 3. Initialize Pager
		pager, writer, err := pathcmds.NewPager(flagPage)
		if err != nil {
			return fmt.Errorf("failed to initialize pager: %w", err)
		}
		defer pager.Close()

		// 4. Format and Print
		formatter := pathcmds.NewFormatter(writer)
		if err := formatter.FormatFolders(filtered); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// CLI Flags
	RootCmd.Flags().BoolVarP(&flagSystem, "system", "s", false, "Filter and show system command sets (/bin, /usr/bin, /sbin, /usr/sbin)")
	RootCmd.Flags().BoolVarP(&flagUser, "user", "u", false, "Filter and show user command sets (.local/bin, homebrew/bin, /usr/local/bin, etc.)")
	RootCmd.Flags().BoolVarP(&flagApps, "apps", "a", false, "Filter and show application-specific folders (.dotnet/sdk, rust/cargo, etc.)")
	RootCmd.Flags().BoolVarP(&flagPage, "page", "p", false, "Page the final output using the system's 'less -R' command")
}
