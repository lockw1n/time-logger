package main

import (
	"log"
	"time"

	"github.com/lockw1n/time-logger/internal/app"
	"github.com/lockw1n/time-logger/internal/migration"
)

func main() {
	log.Println("🔧 Starting DB migration service...")
	database := app.RetryConnect(5, 2*time.Second)

	log.Println("🛠️ Running migrations...")
	if err := database.AutoMigrate(migration.ModelsForMigration()...); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	log.Println("🎉 Migration completed successfully.")
}
