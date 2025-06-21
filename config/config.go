package config

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yaml
var config []byte

type Config struct {
	PasswordHasher struct {
		Cost int `yaml:"cost" validate:"required,min=1"`
	} `yaml:"password_hasher"`

	// TODO: add further configuration e.g. JWT secret, database connection string, etc.
}

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	validator := validator.New()
	if err := validator.Struct(c); err != nil {
		return err
	}
	return nil
}
