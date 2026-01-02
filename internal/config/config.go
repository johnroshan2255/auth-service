package config

import (
	"os"
)

type Config struct {
	DBUrl  string
	Port   string
	JWTKey string
	Redis  string
}

func LoadConfig() *Config {
	return &Config{
		DBUrl:  os.Getenv("POSTGRES_URL"),
		Port:   os.Getenv("PORT"),
		JWTKey: os.Getenv("JWT_KEY"),
		Redis:  os.Getenv("REDIS_URL"),
	}
}
