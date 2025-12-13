package config

import "os"

type Config struct {
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	RedisHost         string
	RedisPassword     string
	RabbitMQURL       string
	ServicePort       string
	ConsulAddr        string
	PlayerServiceURL  string
	PokemonServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "battle_db"),
		DBUser:            getEnv("DB_USER", "maushold"),
		DBPassword:        getEnv("DB_PASSWORD", "changeme"),
		RedisHost:         getEnv("REDIS_HOST", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ServicePort:       getEnv("SERVICE_PORT", "8003"),
		ConsulAddr:        getEnv("CONSUL_ADDR", "consul:8500"),
		PlayerServiceURL:  getEnv("PLAYER_SERVICE_URL", "http://player-service:8001"),
		PokemonServiceURL: getEnv("POKEMON_SERVICE_URL", "http://pokemon-service:8002"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
