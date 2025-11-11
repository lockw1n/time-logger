package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/lockw1n/time-logger/internal/handlers"
	"github.com/lockw1n/time-logger/internal/health"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Health routes
	health.RegisterRoutes(r, db)

	// Entries routes
	entryHandler := handlers.NewEntryHandler(db)
	api := r.Group("/api")
	{
		api.GET("/entries", entryHandler.List)
		api.GET("/entries/:id", entryHandler.Get)
		api.POST("/entries", entryHandler.Create)
		api.PUT("/entries/:id", entryHandler.Update)
		api.DELETE("/entries/:id", entryHandler.Delete)

		// Delete all entries for a ticket (row delete)
		api.DELETE("/tickets/:ticket", entryHandler.DeleteByTicket)
	}

	return r
}
