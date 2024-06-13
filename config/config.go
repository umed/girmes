package config

import (
	"github.com/caarlos0/env/v11"
)

type SecretString string

func (SecretString) String() string {
	return "********"
}

type Config struct {
	DebugMode         bool   `env:"DEBUG_MODE"          envDefault:"false"`
	LogLevel          string `env:"LOG_LEVEL"           envDefault:"debug"`
	GitHubAccessToken string `env:"GITHUB_ACCESS_TOKEN"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
