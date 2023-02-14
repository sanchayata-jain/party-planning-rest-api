package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Host string `envconfig:"HOST" default:"127.0.0.1"`
	Port string `envconfig:"PORT" default:"3000"`
}

func Load() (*Config, error) {
	var c Config
	return &c, envconfig.Process("cf", &c)
}
