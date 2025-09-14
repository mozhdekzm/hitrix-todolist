package config

import (
	"os"
)

type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr   string
	RedisStream string

	ServerPort string
}

func Load() *Config {
	cfg := &Config{
		DBDriver:   getEnv("DB_DRIVER", "mysql"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "todos"),

		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisStream: getEnv("REDIS_STREAM", "todo-stream"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
