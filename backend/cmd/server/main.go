package main

import (
	"log"
	"os"
	"time"

	"github.com/lockw1n/time-logger/internal/db"
	"github.com/lockw1n/time-logger/internal/models"
	"github.com/lockw1n/time-logger/internal/router"
	"gorm.io/gorm"
)

// retryConnect tries to connect to DB with exponential backoff
func retryConnect(maxRetries int, delay time.Duration) *gorm.DB {
	var database *gorm.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		database = db.Connect()
		sqlDB, _ := database.DB()
		if err = sqlDB.Ping(); err == nil {
			log.Printf("✅ Connected to database after %d attempt(s)", i)
			return database
		}
		log.Printf("⚠️ Database not ready (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(delay)
		delay *= 2 // exponential backoff
	}

	log.Fatalf("❌ Could not connect to database after %d attempts: %v", maxRetries, err)
	return nil
}

func main() {
	database := retryConnect(5, 2*time.Second)

	log.Println("🛠️ Running migrations...")
	if err := database.AutoMigrate(&models.Entry{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Database schema ready")

	r := router.SetupRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
