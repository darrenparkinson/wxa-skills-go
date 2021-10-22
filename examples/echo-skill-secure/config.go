package main

import (
	"os"

	"github.com/spf13/viper"
)

// Config holds a reference to any required config.  It uses viper to
// unmarshal environment variables into the configuration.
type Config struct {
	Port  int
	Skill struct {
		PrivateKey string `mapstructure:"skill_private_key"`
		PublicKey  string `mapstructure:"skill_public_key"`
		Secret     string `mapstructure:"skill_secret"`
	} `mapstructure:",squash"`
}

func loadConfig() (*Config, error) {

	// Set defaults
	viper.SetDefault("port", 8080)
	viper.SetDefault("skill_private_key", "")
	viper.SetDefault("skill_public_key", "")
	viper.SetDefault("skill_secret", "")

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

	if cfg.Skill.PrivateKey == "" || cfg.Skill.PublicKey == "" || cfg.Skill.Secret == "" {
		// let's try to find the files instead?
		if cfg.Skill.PrivateKey == "" {
			f, err := os.ReadFile("private.pem")
			if err != nil {
				return nil, ErrMissingEnvironment
			}
			cfg.Skill.PrivateKey = string(f)
		}
		if cfg.Skill.PublicKey == "" {
			f, err := os.ReadFile("public.pem")
			if err != nil {
				return nil, ErrMissingEnvironment
			}
			cfg.Skill.PublicKey = string(f)
		}
		if cfg.Skill.Secret == "" {
			f, err := os.ReadFile("secret.txt")
			if err != nil {
				return nil, ErrMissingEnvironment
			}
			cfg.Skill.Secret = string(f)
		}
	}

	return &cfg, nil
}
