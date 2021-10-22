package command

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/darrenparkinson/wxa-skills-go/pkg/wxaskillsservice"
	"github.com/mitchellh/cli"
)

// DeleteSkillCommand provides the entry point for the command
type DeleteSkillCommand struct {
	UI cli.Ui
}

// Help provies the help text for this command.
func (c *DeleteSkillCommand) Help() string {
	helpText := `
Usage: wxa-cli [global options] delete-skill [options]

  Delete skill on the skills service.

Options:
  -id=ID           The Skill ID
  -hard            Pass the HARD_DELETE flag.
  -token=TOKEN     Your personal access token from developer.webex.com.
  -developerid=ID  Your base64 decoded developer id.
`
	return strings.TrimSpace(helpText)
}

// Run provides the command functionality
func (c *DeleteSkillCommand) Run(args []string) int {
	var id, token, developerID string
	var hard bool
	cmdFlags := flag.NewFlagSet("deleteskill", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&id, "id", "", "the skill id to delete")
	cmdFlags.BoolVar(&hard, "hard", false, "pass the HARD_DELETE flag")
	cmdFlags.StringVar(&token, "token", "", "your personal access token")
	cmdFlags.StringVar(&developerID, "developerid", "", "your base64 decoded developer id")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	if id == "" || token == "" || developerID == "" {
		c.UI.Error("error: missing required flags")
		return 1
	}
	ctx := context.Background()
	ss, err := wxaskillsservice.NewClient(developerID, token, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	err = ss.DeleteSkill(ctx, id, hard)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.UI.Output(fmt.Sprintf("Skill %s deleted.", id))

	return 0
}

// Synopsis provides the one liner
func (c *DeleteSkillCommand) Synopsis() string {
	return "Delete skill on the skills service."
}
