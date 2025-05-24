package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Port: getEnvOrDefault("PORT", "3333"),
	}, nil
}

func getEnvOrDefault(envVariable string, defaultValue string) string {
	var env = os.Getenv(envVariable)
	if env == "" {
		return defaultValue
	}
	return env
}
