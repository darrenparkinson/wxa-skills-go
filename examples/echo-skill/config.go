package main

import (
	"github.com/spf13/viper"
)

// Config holds a reference to any required config
type Config struct {
	Port int
}

func loadConfig() (*Config, error) {

	// Set defaults
	viper.SetDefault("port", 8080)

	// Set up Viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.ReadInConfig() // ignore error since there may not be a config file
	viper.AutomaticEnv() // read in environment variables that match

	cfg := Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
