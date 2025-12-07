package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg PostgresConfig) (*sql.DB, error) {
	// Format connection string
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Open a connection
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Unable to ping database: %v\n", err)
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully!")
	return db, nil
}
