package config

import "github.com/PatrickKvartsaniy/print-me-at/repository"

type Config struct {
	Port            int
	Redis           repository.RedisConfig
	HealthCheckPort int
}

func ReadOS() Config {
	return Config{}
}
