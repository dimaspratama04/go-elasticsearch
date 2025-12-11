package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresSSLMode  string
	ElasticURL       string
	ElasticUsername  string
	ElasticPassword  string
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUsername string
	RabbitMQPassword string
	RabbitMQVhost    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		PostgresHost:     getEnv("PG_HOST", "localhost"),
		PostgresPort:     getEnv("PG_PORT", "5432"),
		PostgresUser:     getEnv("PG_USERNAME", "postgres"),
		PostgresPassword: getEnv("PG_PASSWORD", "password"),
		PostgresDBName:   getEnv("PG_DATABASE", "postgres"),
		PostgresSSLMode:  getEnv("PG_SSLMODE", "disable"),
		ElasticURL:       getEnv("ELASTIC_URL", "http://localhost:9200"),
		ElasticUsername:  getEnv("ELASTIC_USERNAME", "elastic"),
		ElasticPassword:  getEnv("ELASTIC_PASSWORD", "changeme"),
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUsername: getEnv("RABBITMQ_USERNAME", "root"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "changeme"),
		RabbitMQVhost:    getEnv("RABBITMQ_VHOST", "/"),
	}

	log.Println("Configuration loaded successfully")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
