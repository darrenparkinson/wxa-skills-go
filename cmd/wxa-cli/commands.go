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
