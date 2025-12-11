package config

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgresConnection(host, username, password, dbname, port, sslmode string) *gorm.DB {
	// Format connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, username, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}

	// Get underlying sql.DB untuk connection pool tuning
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed get sql.DB: %v", err)
	}

	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL")

	return db
}
