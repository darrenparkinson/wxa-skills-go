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
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/darrenparkinson/wxa-skills-go/pkg/wxaskillsservice"
	"github.com/mitchellh/cli"
)

// ListSkillsCommand provides the entry point for the command
type ListSkillsCommand struct {
	UI cli.Ui
}

// Help provies the help text for this command.
func (c *ListSkillsCommand) Help() string {
	helpText := `
Usage: wxa-cli [global options] list-skills [options]

  List skills configured on the skills service.

Options:
  -token=TOKEN     Your personal access token from developer.webex.com.
  -developerid=ID  Your base64 decoded developer id.

`
	return strings.TrimSpace(helpText)
}

// Run provides the command functionality
func (c *ListSkillsCommand) Run(args []string) int {
	var token, developerID string
	cmdFlags := flag.NewFlagSet("listskills", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&token, "token", "", "your personal access token")
	cmdFlags.StringVar(&developerID, "developerid", "", "your base64 decoded developer id")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	if token == "" || developerID == "" {
		c.UI.Error("error: token and developer id flags reqUIred")
		return 1
	}
	ctx := context.Background()
	ss, err := wxaskillsservice.NewClient(developerID, token, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	skills, err := ss.ListSkills(ctx)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	for _, s := range skills {
		var deleted string
		if *s.Deleted {
			deleted = "(SOFT DELETED)"
		}
		c.UI.Output(fmt.Sprintf("ID: %s; DeveloperID: %s; URL: %s; Name: %s; ContactEmail: %s; %s", *s.SkillID, *s.DeveloperID, *s.URL, *s.Name, *s.ContactEmail, deleted))
	}

	return 0
}

// Synopsis provides the one liner
func (c *ListSkillsCommand) Synopsis() string {
	return "List skills configured on the skills service."
}
