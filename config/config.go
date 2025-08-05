package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	ServerPort   string
	AWSRegion    string
	AWSEndpoint  string
	SQSQueueName string
}

func LoadConfig() Config {
	cfg := Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       "5432",
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "todo"),
		DBName:       getEnv("DB_NAME", "todo"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
		AWSEndpoint:  getEnv("AWS_ENDPOINT", "http://localstack:4566"),
		SQSQueueName: getEnv("SQS_QUEUE_NAME", "todo-queue"),
	}
	fmt.Println(cfg)
	return cfg
}

func (c Config) SQSQueueURL() string {
	return fmt.Sprintf("%s/000000000000/%s", c.AWSEndpoint, c.SQSQueueName)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
