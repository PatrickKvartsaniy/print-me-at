package config

import "github.com/PatrickKvartsaniy/print-me-at/repository"

type Config struct {
	Port  int
	Redis repository.RedisConfig
}

func ReadOS() Config {
	return Config{}
}
