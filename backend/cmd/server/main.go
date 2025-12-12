package main

import (
	"log"
	"os"
	"time"

	"github.com/lockw1n/time-logger/internal/app"
	"github.com/lockw1n/time-logger/internal/router"
)

func main() {
	log.Println("🌐 Starting API server...")

	database := app.RetryConnect(5, 2*time.Second)
	container := app.NewContainer(database)

	r := router.SetupRouter(container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API running at :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
