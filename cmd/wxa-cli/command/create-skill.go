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

// CreateSkillCommand provides the entry point for the command
type CreateSkillCommand struct {
	UI cli.Ui
}

// Help provies the help text for this command.
func (c *CreateSkillCommand) Help() string {
	helpText := `
Usage: wxa-cli [global options] create-skill [options]

  Create skill on the skills service.

Options:
  -name=NAME       The name of your skill.
  -url=URL         The publicly accessible url for your skill.
  -contact=EMAIL   The contact email address for the skill.
  -token=TOKEN     Your personal access token from developer.webex.com.
  -developerid=ID  Your base64 decoded developer id.
  -public=KEY      The public key for your skill.
  -secret=SECRET   The secret for your skill.
`
	return strings.TrimSpace(helpText)
}

// Run provides the command functionality
func (c *CreateSkillCommand) Run(args []string) int {
	var name, url, contact, public, secret, token, developerID string
	cmdFlags := flag.NewFlagSet("listskills", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&name, "name", "", "the name of your skill.")
	cmdFlags.StringVar(&url, "url", "", "the publicly accessible url for your skill.")
	cmdFlags.StringVar(&contact, "contact", "", "the contact email address for the skill.")
	cmdFlags.StringVar(&public, "public", "", "the public key for your skill")
	cmdFlags.StringVar(&secret, "secret", "", "the secret for your skill")
	cmdFlags.StringVar(&token, "token", "", "your personal access token")
	cmdFlags.StringVar(&developerID, "developerid", "", "your base64 decoded developer id")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	if name == "" || url == "" || contact == "" || token == "" || developerID == "" || public == "" || secret == "" {
		c.UI.Error("error: missing required flags")
		return 1
	}
	ctx := context.Background()
	ss, err := wxaskillsservice.NewClient(developerID, token, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	newSkill := wxaskillsservice.Skill{
		Name:         wxaskillsservice.String(name),
		URL:          wxaskillsservice.String(url),
		ContactEmail: wxaskillsservice.String(contact),
		PublicKey:    wxaskillsservice.String(public),
		Secret:       wxaskillsservice.String(secret),
		Languages:    []string{"en"},
	}
	skill, err := ss.CreateSkill(ctx, newSkill)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.UI.Output(fmt.Sprintf("Skill %s created with URL %s (ID:%s)", *skill.Name, *skill.URL, *skill.SkillID))

	return 0
}

// Synopsis provides the one liner
func (c *CreateSkillCommand) Synopsis() string {
	return "List skills configured on the skills service."
}
