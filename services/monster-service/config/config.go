package config

import (
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	RedisHost     string
	RedisPassword string
	RabbitMQURL   string
	ServicePort   string
	ConsulAddr    string
}

func LoadConfig() *Config {
	c := &Config{
		DBHost:        getEnv("DB_HOST"),
		DBPort:        getEnv("DB_PORT"),
		DBName:        getEnv("DB_NAME"),
		DBUser:        getEnv("DB_USER"),
		DBPassword:    getEnv("DB_PASSWORD"),
		RedisHost:     getEnv("REDIS_HOST"),
		RedisPassword: getEnv("REDIS_PASSWORD"),
		RabbitMQURL:   getEnv("RABBITMQ_URL"),
		ServicePort:   getEnv("SERVICE_PORT"),
		ConsulAddr:    getEnv("CONSUL_ADDR"),
	}
	return c
}

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return ""
}
