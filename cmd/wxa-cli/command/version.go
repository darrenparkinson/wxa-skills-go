package command

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/cli"
)

// VersionCommand is the top level struct for the cli VersionCommand.
// It holds a reference to the cli.Ui for logging etc.
type VersionCommand struct {
	Version string
	UI      cli.Ui
}

// Help provies the help text for this command.
func (c *VersionCommand) Help() string {
	return ""
}

// Run provides the command functionality
func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "wxa-cli v%s", c.Version)
	c.UI.Output(versionString.String())
	return 0
}

// Synopsis provides the one liner
func (c *VersionCommand) Synopsis() string {
	return "Show version information."
}
