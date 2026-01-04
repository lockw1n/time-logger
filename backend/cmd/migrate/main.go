package main

import (
	"log"
	"time"

	"github.com/lockw1n/time-logger/internal/app"
	"github.com/lockw1n/time-logger/internal/migration"
)

func main() {
	log.Println("🔧 Starting DB migration service...")
	db := app.RetryConnect(5, 2*time.Second)

	log.Println("🛠️ Running structural migrations (AutoMigrate)...")
	if err := db.AutoMigrate(migration.ModelsForMigration()...); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	log.Println("🛠️ Running explicit migrations...")
	for _, m := range migration.ExplicitMigrations() {
		log.Printf("➡️ %s", m.Name())
		if err := m.Run(db); err != nil {
			log.Fatalf("❌ Migration failed (%s): %v", m.Name(), err)
		}
	}

	log.Println("🎉 Migration completed successfully.")
}
