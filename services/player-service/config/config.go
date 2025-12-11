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
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "player_db"),
		DBUser:        getEnv("DB_USER", "maushold"),
		DBPassword:    getEnv("DB_PASSWORD", "changeme"),
		RedisHost:     getEnv("REDIS_HOST", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ServicePort:   getEnv("SERVICE_PORT", "8001"),
		ConsulAddr:    getEnv("CONSUL_ADDR", "consul:8500"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
