package command

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

// GenerateKeysCommand provides the entry point for the command
type GenerateKeysCommand struct {
	UI cli.Ui
}

// Help provies the help text for this command.
func (c *GenerateKeysCommand) Help() string {
	helpText := `
Usage: wxa-cli [global options] generate-keys [options]

  Generate an RSA keypair in pem format.

Options:
  -public=FILENAME    Specify the filename for the generated public key. Default "public.pem".

  -private=FILENAME   Specify the filename for the generated private key. Default "private.pem".

`
	return strings.TrimSpace(helpText)
}

// Run provides the command functionality
func (c *GenerateKeysCommand) Run(args []string) int {
	var privateFilename, publicFilename string
	cmdFlags := flag.NewFlagSet("generatekeys", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&publicFilename, "public", "public.pem", "public key file to create")
	cmdFlags.StringVar(&privateFilename, "private", "private.pem", "private key file to create")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	err := generateKeys(privateFilename, publicFilename)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return 0
}

// Synopsis provides the one liner
func (c *GenerateKeysCommand) Synopsis() string {
	return "Generate an RSA keypair in pem format."
}

func generateKeys(private, public string) error {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("cannot generate RSA key: %s", err)
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(private)
	if err != nil {
		return fmt.Errorf("error creating private key file: %s", err)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("error encoding private pem: %s", err)
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return fmt.Errorf("error creating public key: %s", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create(public)
	if err != nil {
		return fmt.Errorf("error creating public key file: %s", err)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		return fmt.Errorf("error encoding public pem: %s", err)
	}
	return nil
}
