package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/app"
	"github.com/lockw1n/time-logger/internal/health"

	activityhandler "github.com/lockw1n/time-logger/internal/activity/handler"
	companyhandler "github.com/lockw1n/time-logger/internal/company/handler"
	consultanthandler "github.com/lockw1n/time-logger/internal/consultant/handler"
	contracthandler "github.com/lockw1n/time-logger/internal/contract/handler"
	entryhandler "github.com/lockw1n/time-logger/internal/entry/handler"
	invoicehandler "github.com/lockw1n/time-logger/internal/invoice/handler"
	tickethandler "github.com/lockw1n/time-logger/internal/ticket/handler"
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
	contractHandler := contracthandler.NewHandler(container.ContractService)
	ticketHandler := tickethandler.NewHandler(container.TicketService)
	activityHandler := activityhandler.NewHandler(container.ActivityService)
	entryHandler := entryhandler.NewHandler(container.EntryService)
	timesheetHandler := timesheethandler.NewHandler(container.TimesheetService)
	invoiceHandler := invoicehandler.NewHandler(container.InvoiceService, container.PdfRenderer)

	api := r.Group("/api")
	{
		api.POST("/companies", companyHandler.CreateCompany)
		api.PUT("/companies/:id", companyHandler.UpdateCompany)
		api.DELETE("/companies/:id", companyHandler.DeleteCompany)
		api.GET("/companies/:id", companyHandler.GetCompany)
		api.GET("/companies", companyHandler.ListCompanies)

		api.POST("/consultants", consultantHandler.CreateConsultant)
		api.PUT("/consultants/:id", consultantHandler.UpdateConsultant)
		api.DELETE("/consultants/:id", consultantHandler.DeleteConsultant)
		api.GET("/consultants/:id", consultantHandler.GetConsultant)

		api.POST("/contracts", contractHandler.CreateContract)
		api.PUT("/contracts/:id", contractHandler.UpdateContract)
		api.DELETE("/contracts/:id", contractHandler.DeleteContract)
		api.GET("/contracts/:id", contractHandler.GetContract)
		api.GET("/contracts", contractHandler.ListContractsForConsultant)

		api.POST("/tickets", ticketHandler.CreateTicket)
		api.PUT("/tickets/:id", ticketHandler.UpdateTicket)
		api.DELETE("/tickets/:id", ticketHandler.DeleteTicket)
		api.GET("/tickets/:id", ticketHandler.GetTicket)
		api.GET("/tickets", ticketHandler.ListTicketsForCompany)

		api.POST("/activities", activityHandler.CreateActivity)
		api.PUT("/activities/:id", activityHandler.UpdateActivity)
		api.DELETE("/activities/:id", activityHandler.DeleteActivity)
		api.GET("/activities/:id", activityHandler.GetActivity)
		api.GET("/activities", activityHandler.ListActivitiesForCompany)

		api.POST("/entries", entryHandler.CreateEntry)
		api.PUT("/entries/:id", entryHandler.UpdateEntry)
		api.DELETE("/entries/:id", entryHandler.DeleteEntry)
		api.GET("/entries/:id", entryHandler.GetEntry)

		api.GET("/timesheet", timesheetHandler.GetTimesheet)

		api.POST("/invoices/generate", invoiceHandler.GenerateInvoice)
	}

	return r
}
