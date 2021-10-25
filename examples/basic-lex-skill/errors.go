package main

import (
	"errors"
)

var (
	// ErrMissingEnvironment is a constant error for missing environment variables
	ErrMissingEnvironment = errors.New("missing required environment variables or files for keys or secret")
	// ErrMissingAWSEnvironment is a constant error for missing aws environment variables
	ErrMissingAWSEnvironment = errors.New("missing required environment variables for aws or lex")
	// ErrMissingWeatherEnvironment is a constant error for missing weather environment variables
	ErrMissingWeatherEnvironment = errors.New("missing required environment variables for weather service")
)
