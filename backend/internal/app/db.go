package app

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// RetryConnect tries to connect to DB with exponential backoff.
func RetryConnect(maxRetries int, delay time.Duration) *gorm.DB {
	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db = Connect()
		if db != nil {
			sqlDB, _ := db.DB()
			if err = sqlDB.Ping(); err == nil {
				log.Printf("✅ Connected to database on attempt %d", attempt)
				return db
			}
		}

		log.Printf("⚠️ Database not ready (attempt %d/%d): %v", attempt, maxRetries, err)
		time.Sleep(delay)
		delay *= 2
	}

	log.Fatalf("❌ Could not connect to database after %d attempts: %v", maxRetries, err)
	return nil
}

// Connect initializes the database connection.
func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ DB connection failed: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("❌ DB connection pool error: %v", err)
		return nil
	}

	// Recommended pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
