package pathcmds

import (
	"testing"
)

func TestCategorize(t *testing.T) {
	c := NewCategorizer()

	tests := []struct {
		path     string
		expected string
	}{
		// System paths
		{"/bin", CategorySystem},
		{"/usr/bin", CategorySystem},
		{"/sbin", CategorySystem},
		{"/usr/sbin", CategorySystem},
		{"/bin/subdir", CategorySystem},

		// User paths
		{"/usr/local/bin", CategoryUser},
		{"/usr/local/sbin", CategoryUser},
		{"/opt/homebrew/bin", CategoryUser},
		{"/opt/homebrew/sbin", CategoryUser},
		{"/Users/john/bin", CategoryUser},
		{"/Users/john/.local/bin", CategoryUser},
		{"/home/linuxuser/bin", CategoryUser},

		// App paths
		{"/Users/john/.cargo/bin", CategoryApp},
		{"/Users/john/.rustup/toolchains/bin", CategoryApp},
		{"/Users/john/.dotnet", CategoryApp},
		{"/Users/john/.dotnet/tools", CategoryApp},
		{"/Users/john/go/bin", CategoryApp},
		{"/Users/john/.nvm/versions/node/v16.0.0/bin", CategoryApp},
		{"/Users/john/.n/bin", CategoryApp},
		{"/Applications/Obsidian.app/Contents/MacOS", CategoryApp},
		{"/opt/aws/bin", CategoryApp},
		{"/opt/microsoft/powershell/bin", CategoryApp},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			actual := c.Categorize(tt.path)
			if actual != tt.expected {
				t.Errorf("Categorize(%q) = %q; want %q", tt.path, actual, tt.expected)
			}
		})
	}
}
