package main

import (
	"os"

	"github.com/darrenparkinson/wxa-skills-go/cmd/wxa-cli/command"
	"github.com/mitchellh/cli"
)

// Commands holds a map of each of the cli commands
var Commands map[string]cli.CommandFactory

func init() {
	bui := &cli.BasicUi{Writer: os.Stdout}
	ui := &cli.ColoredUi{
		Ui:          bui,
		OutputColor: cli.UiColorBlue,
		InfoColor:   cli.UiColorGreen,
		WarnColor:   cli.UiColorYellow,
		ErrorColor:  cli.UiColorRed,
	}
	Commands = map[string]cli.CommandFactory{
		"generate-keys": func() (cli.Command, error) {
			return &command.GenerateKeysCommand{UI: ui}, nil
		},
		"generate-secret": func() (cli.Command, error) {
			return &command.GenerateSecretCommand{UI: ui}, nil
		},
		"list-skills": func() (cli.Command, error) {
			return &command.ListSkillsCommand{UI: ui}, nil
		},
		"create-skill": func() (cli.Command, error) {
			return &command.CreateSkillCommand{UI: ui}, nil
		},
		"delete-skill": func() (cli.Command, error) {
			return &command.DeleteSkillCommand{UI: ui}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				UI:      ui,
			}, nil
		},
	}
}
