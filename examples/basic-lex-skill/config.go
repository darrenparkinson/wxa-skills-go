package main

import (
	"os"

	"github.com/spf13/viper"
)

// Config holds a reference to any required config.  It uses viper to
// unmarshal environment variables into the configuration.
type Config struct {
	Port           int
	OpenWeatherMap struct {
		APIKey string `mapstructure:"openweathermap_apikey"`
	} `mapstructure:",squash"`
	Skill struct {
		PrivateKey string `mapstructure:"skill_private_key"`
		PublicKey  string `mapstructure:"skill_public_key"`
		Secret     string `mapstructure:"skill_secret"`
	} `mapstructure:",squash"`
	AWS struct {
		Region          string `mapstructure:"aws_region"`
		AccessKeyID     string `mapstructure:"aws_access_key_id"`
		SecretAccessKey string `mapstructure:"aws_secret_access_key"`
	} `mapstructure:",squash"`
	Lex struct {
		Alias   string `mapstructure:"lex_alias"`
		BotName string `mapstructure:"lex_botname"`
	} `mapstructure:",squash"`
}

func loadConfig() (*Config, error) {

	// Set defaults
	viper.SetDefault("port", 8080)
	viper.SetDefault("skill_private_key", "")
	viper.SetDefault("skill_public_key", "")
	viper.SetDefault("skill_secret", "")

	viper.SetDefault("aws_region", "")
	viper.SetDefault("aws_access_key_id", "")
	viper.SetDefault("aws_secret_access_key", "")
	viper.SetDefault("lex_alias", "")
	viper.SetDefault("lex_botname", "")
	viper.SetDefault("openweathermap_apikey", "")

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

	if cfg.OpenWeatherMap.APIKey == "" {
		return nil, ErrMissingWeatherEnvironment
	}

	if cfg.AWS.Region == "" || cfg.AWS.AccessKeyID == "" || cfg.AWS.SecretAccessKey == "" || cfg.Lex.Alias == "" || cfg.Lex.BotName == "" {
		return nil, ErrMissingAWSEnvironment
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
