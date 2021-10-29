# Webex Assistant Skills - Go

[![Status](https://img.shields.io/badge/status-wip-yellow)](https://github.com/darrenparkinson/wxa-skills-go) ![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/darrenparkinson/wxa-skills-go) ![GitHub](https://img.shields.io/github/license/darrenparkinson/wxa-skills-go?color=brightgreen) [![GoDoc](https://pkg.go.dev/badge/darrenparkinson/wxa-skills-go)](https://pkg.go.dev/github.com/darrenparkinson/wxa-skills-go) [![Go Report Card](https://goreportcard.com/badge/github.com/darrenparkinson/wxa-skills-go)](https://goreportcard.com/report/github.com/darrenparkinson/wxa-skills-go)

This repository holds example skills and a cli utility written in Go for interacting with Webex Assistant Skills.  For more information on Webex Assistant Skills, please see the official [Webex Assistant Skills Overview](https://developer.webex.com/docs/webex-assistant-skills-overview).

It is intended as an alternative to the official [Webex Assistant SDK](https://github.com/cisco/webex-assistant-sdk), as documented in the [Webex Assistant Skills Guide](https://developer.webex.com/docs/api/guides/webex-assistant-skills-guide), without the need for installing python and the required dependencies.  Simply download the binaries from the releases tab to get started.  

It is currently a work in progress, so expect things to change as more is learned about this new webex feature.

## Quick Start

You should be famililar with the [Getting Started](https://github.com/cisco/webex-assistant-sdk/tree/main/get_started_documentation#create-skill-on-skills-service) documentation. This quickstart assumes you'll be using the simulator to test.

To get started, you will need:

* [ ] Your tenant enabled for skills;
* [ ] Your personal access token to register the skill: [get your token here](https://developer.webex.com/docs/api/v1/people/list-people) from the Copy button in the Authorization section;
* [ ] Your base64 decoded developer ID: [get your ID](https://developer.webex.com/docs/api/v1/people/get-my-own-details) and [base64 decode it here](https://www.base64decode.org/)  taking the last part after `ciscospark://us/PEOPLE/`
* [ ] Your base64 decoded organisation ID (for the simulator): [get your orgId](https://developer.webex.com/docs/api/v1/people/get-my-own-details) and [base64 decode it here](https://www.base64decode.org/) taking the last part after `ciscospark://us/ORGANIZATION/`
* [ ] A token with the `assistant` scope to run the skill.  You can temporarily [get this from here](https://3bfnei7xs2.execute-api.us-east-1.amazonaws.com/production/wxa-token) until there is proper tooling;


1. Create a folder and download the binaries from the [releases page](https://github.com/darrenparkinson/wxa-skills-go/releases):  
    a. `wxa-cli` - for generating the keys and the secret, along with registering your skill;  
    b. `echo-skill-secure` - for the test skill;  
    c. `echo-skill-secure-tester` for testing the skill locally;  

2. Generate a public/private key pair: 
```sh
$ ./wxa-cli generate-keys
```

3. Generate a secret: 
```sh
$ ./wxa-cli generate-secret > secret.txt
```

4. (Optional) Set your environment variables or put them in a .env file:
```sh
SKILL_PUBLIC_KEY=<YOUR PUBLIC KEY HERE>
SKILL_PRIVATE_KEY=<YOUR PRIVATE KEY HERE>
SKILL_SECRET=<YOUR SECRET HERE>
```

If you don't set these, by default the skill will look in the current directory for `secret.txt`, `private.pem` and `public.pem`.

5. Set up a tunnel to your machine using localtunnel or [ngrok](https://ngrok.com), e.g:

```sh
ngrok http 8080
```

Use the `https` endpoint provided by ngrok in the next step.

6. Create the Skill on the Skills Service using the details obtained earlier:

You can do this on the [Webex Assistant Skills Developer Portal](https://skills-developer.intelligence.webex.com/) as documented in the [Webex Assistant Skills Guide Developer Portal Guide](https://developer.webex.com/docs/api/guides/webex-assistant-skills-guide-developer-portal-guide) or use the `wxa-cli` command:

```sh
$ wxa-cli create-skill --name="Echo" --url="<YOUR_URL_FROM_STEP_5>" --contact="<YOUR_EMAIL>" -secret="$(cat secret.txt)" --public="$(cat public.pem)" --token="<YOUR_PERSONAL_ACCESS_TOKEN>" --developerid="<YOUR_DEVELOPER_ID>"
```

Replacing values in `< >` with the relevant details from earlier. Note the use of `cat` to provide the secret.txt and public.pem content into the command.

7. Run the skill:
```sh
$ ./echo-skill-secure
```

8. Test the skill locally:
```sh
$ ./echo-skill-secure-tester
```

9. Test the skill with the simulator:

* Visit `https://assistant-web.intelligence.webex.com/`.  
* Enter the **base64 decoded** organisation ID and the assistant scoped token.  
* Say or type `ask echo hello there`

# Installation

## Binaries

The simplest way is to download the binaries from the [releases page](https://github.com/darrenparkinson/wxa-skills-go/releases).

## Go Install

If you have Go installed, you can install using the following commands:

```sh
$ go install github.com/darrenparkinson/wxa-skills-go/cmd/wxa-cli@latest
$ go install github.com/darrenparkinson/wxa-skills-go/examples/echo-skill-secure@latest
$ go install github.com/darrenparkinson/wxa-skills-go/examples/echo-skill-secure/cmd/echo-skill-secure-tester@latest
```

## Compiling from source

Again, if you have Go installed, you can also compile from source:

```sh
$ git clone https://github.com/darrenparkinson/wxa-skills-go
$ cd wxa-skills-go
$ go mod tidy
$ go build ./cmd/wxa-cli
$ go build ./examples/echo-skill-secure
$ go build ./examples/echo-skill-secure/echo-skill-secure-tester
$ ./wxa-cli --version
```