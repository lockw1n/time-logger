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
	entryService := handlers.NewEntryService(db)
	entryHandler := handlers.NewEntryHandler(entryService)
	invoiceHandler := handlers.NewInvoiceHandler(entryService)
	api := r.Group("/api")
	{
		api.GET("/entries", entryHandler.List)
		api.GET("/entries/:id", entryHandler.Get)
		api.POST("/entries", entryHandler.Create)
		api.PUT("/entries/:id", entryHandler.Update)
		api.DELETE("/entries/:id", entryHandler.Delete)

		api.GET("/tickets/summary", entryHandler.Summary)
		// Delete all entries for a ticket (row delete)
		api.DELETE("/tickets/:ticket", entryHandler.DeleteByTicket)

		api.GET("/reports/monthly", entryHandler.MonthlyReport)
		api.POST("/reports/invoice/pdf", invoiceHandler.GeneratePDF)
	}

	return r
}
