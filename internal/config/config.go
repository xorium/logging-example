package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	// https://pkg.go.dev/github.com/rs/zerolog#LevelFieldName
	LogLevel       string `envconfig:"LOG_LEVEL" default:"trace"`
	HttpListenAddr string `envconfig:"HTTP_LISTEN_ADDR" default:"0.0.0.0:8000"`
	Debug          bool   `envconfig:"DEBUG" default:"false"`
}

func ParseFromEnv() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "can't parse config")
	}

	return cfg, nil
}
