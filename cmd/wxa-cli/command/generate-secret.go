package command

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// GenerateSecretCommand provides the entry point for the command
type GenerateSecretCommand struct {
	UI cli.Ui
}

// Help provies the help text for this command.
func (c *GenerateSecretCommand) Help() string {
	helpText := `
Usage: wxa-cli [global options] generate-secret [options]

  Generate a secret token for signing requests.

`
	return strings.TrimSpace(helpText)
}

// Run provides the command functionality
func (c *GenerateSecretCommand) Run(args []string) int {
	token, err := generateToken()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	// c.UI.Output(token)
	// we print with no newline in order that it doesn't cause a problem when redirected to a file
	fmt.Print(token)
	return 0
}

// Synopsis provides the one liner
func (c *GenerateSecretCommand) Synopsis() string {
	return "Generate a secret token for signing requests."
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	se := base64.URLEncoding.EncodeToString([]byte(b))
	return strings.TrimRight(se, "="), nil
}
