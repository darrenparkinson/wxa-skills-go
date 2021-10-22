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
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

// Version holds the current release version.
const Version = "0.1.0"

func main() {
	// log.SetOutput(ioutil.Discard)
	c := cli.NewCLI("wxa-cli", Version)
	c.Args = os.Args[1:]

	// This just runs the version command instead of the default version.
	for _, arg := range c.Args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(c.Args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], c.Args)
			c.Args = newArgs
			break
		}
	}

	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(exitStatus)
}
