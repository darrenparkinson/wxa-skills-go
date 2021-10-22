package main

import (
	"errors"
)

var (
	// ErrMissingEnvironment is a constant error for missing environment variables
	ErrMissingEnvironment = errors.New("missing required environment variables or files")
)

type webexAssistantModel struct{}

type models struct {
	webexAssistant webexAssistantModel
}

func newModels() models {
	return models{
		webexAssistant: webexAssistantModel{},
	}
}
