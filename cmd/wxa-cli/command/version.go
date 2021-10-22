// Copyright 2021 Darren Parkinson

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
