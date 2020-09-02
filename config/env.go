package config

import (
	"github.com/PatrickKvartsaniy/print-me-at/repository"
	"github.com/spf13/viper"
)

func ReadOS() Config {
	viper.SetDefault("PRETTY_LOG_OUTPUT", true)
	viper.SetDefault("LOG_LEVEL", "DEBUG")

	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("SERVER_PORT", 8080)

	viper.SetDefault("REDIS_POLLING_INTERVAL", "1s")
	viper.SetDefault("REDIS_KEY", "scheduled_messages")
	viper.SetDefault("REDIS_ADDRESS", "0.0.0.0:6379")

	return Config{
		Port: viper.GetInt("SERVER_PORT"),
		Redis: repository.RedisConfig{
			Addr:            viper.GetString("REDIS_ADDRESS"),
			Key:             viper.GetString("REDIS_KEY"),
			PollingInterval: viper.GetDuration("REDIS_POLLING_INTERVAL"),
		},
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		PrettyLogOutput: viper.GetBool("PRETTY_LOG_OUTPUT"),
		LogLevel:        viper.GetString("LOG_LEVEL"),
	}
}
