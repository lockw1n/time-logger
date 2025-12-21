package router

import (
	"github.com/gin-gonic/gin"

	"github.com/lockw1n/time-logger/internal/app"
	companyhandler "github.com/lockw1n/time-logger/internal/handlers/company"
	consultanthandler "github.com/lockw1n/time-logger/internal/handlers/consultant"
	entryhandler "github.com/lockw1n/time-logger/internal/handlers/entry"
	"github.com/lockw1n/time-logger/internal/health"
	invoicehandler "github.com/lockw1n/time-logger/internal/invoice/handler"
	timesheethandler "github.com/lockw1n/time-logger/internal/timesheet/handler"
)

func SetupRouter(container *app.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Health endpoints
	healthHandler := health.NewHandler(container.DB)
	r.GET("/health", healthHandler.Check)
	r.HEAD("/health", healthHandler.Check)

	companyHandler := companyhandler.NewHandler(container.CompanyService)
	consultantHandler := consultanthandler.NewHandler(container.ConsultantService)
	timesheetHandler := timesheethandler.NewTimesheet(container.TimesheetService)
	entryHandler := entryhandler.NewEntryHandler(container.EntryService)
	invoiceHandler := invoicehandler.NewInvoice(container.InvoiceGenerator, container.PdfRenderer, container.ExcelRenderer)

	api := r.Group("/api")
	{
		api.GET("/company", companyHandler.GetCompany)
		api.PUT("/company", companyHandler.UpsertCompany)

		api.GET("/consultant", consultantHandler.GetConsultant)
		api.PUT("/consultant", consultantHandler.UpsertConsultant)

		api.GET("/timesheet", timesheetHandler.GetTimesheet)

		api.POST("/entries", entryHandler.Create)
		api.PUT("/entries/:id", entryHandler.Update)
		api.DELETE("/entries/:id", entryHandler.Delete)

		api.POST("/invoices/monthly", invoiceHandler.GenerateMonthly)
	}

	return r
}
