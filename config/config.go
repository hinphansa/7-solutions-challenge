package config

import (
	"github.com/hinphansa/7-solutions-challenge/pkg/utils"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yaml
var config []byte

type Config struct {
	PasswordHasher struct {
		Cost int `yaml:"cost" validate:"required,min=1"`
	} `yaml:"password_hasher"`

	JWT struct {
		Secret string `yaml:"secret" validate:"required,min=1"`
		TTL    int    `yaml:"ttl" validate:"required,min=1"` // time to live in seconds
	} `yaml:"jwt"`

	Mongo struct {
		URI string `yaml:"uri" validate:"required,min=1"`
		DB  string `yaml:"db" validate:"required,min=1"`
	} `yaml:"mongo"`

	Server struct {
		Port int `yaml:"port" validate:"required,min=1"`
	} `yaml:"server"`

	// TODO: add further configuration e.g. JWT secret, database connection string, etc.
}

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	// validate config with panic, if validation fails
	utils.MustValid(&cfg)
	return &cfg, nil
}
