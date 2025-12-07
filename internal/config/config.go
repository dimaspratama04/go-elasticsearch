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
