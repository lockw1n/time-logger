package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lockw1n/time-logger/internal/app"
	"github.com/lockw1n/time-logger/internal/health"

	activityhandler "github.com/lockw1n/time-logger/internal/activity/handler"
	authhandler "github.com/lockw1n/time-logger/internal/auth/handler"
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

	healthHandler := health.NewHandler(container.DB)
	authHandler := authhandler.NewHandler(container.AuthService)
	companyHandler := companyhandler.NewHandler(container.CompanyService)
	consultantHandler := consultanthandler.NewHandler(container.ConsultantService)
	contractHandler := contracthandler.NewHandler(container.ContractService)
	ticketHandler := tickethandler.NewHandler(container.TicketService)
	activityHandler := activityhandler.NewHandler(container.ActivityService)
	entryHandler := entryhandler.NewHandler(container.EntryService)
	timesheetHandler := timesheethandler.NewHandler(container.TimesheetService)
	invoiceHandler := invoicehandler.NewHandler(container.InvoiceService, container.PdfRenderer)

	r.GET("/health", healthHandler.Check)
	r.HEAD("/health", healthHandler.Check)

	api := r.Group("/api")

	{
		api.POST("/auth/login", authHandler.Login)
	}

	protected := api.Group("")
	protected.Use(container.AuthMiddleware.RequireAuth())
	{
		protected.POST("/companies", companyHandler.CreateCompany)
		protected.PUT("/companies/:id", companyHandler.UpdateCompany)
		protected.DELETE("/companies/:id", companyHandler.DeleteCompany)
		protected.GET("/companies/:id", companyHandler.GetCompany)
		protected.GET("/companies", companyHandler.ListCompanies)

		protected.POST("/consultants", consultantHandler.CreateConsultant)
		protected.PUT("/consultants/:id", consultantHandler.UpdateConsultant)
		protected.DELETE("/consultants/:id", consultantHandler.DeleteConsultant)
		protected.GET("/consultants/:id", consultantHandler.GetConsultant)

		protected.POST("/contracts", contractHandler.CreateContract)
		protected.PUT("/contracts/:id", contractHandler.UpdateContract)
		protected.DELETE("/contracts/:id", contractHandler.DeleteContract)
		protected.GET("/contracts/:id", contractHandler.GetContract)
		protected.GET("/contracts", contractHandler.ListContractsForConsultant)

		protected.POST("/tickets", ticketHandler.CreateTicket)
		protected.PUT("/tickets/:id", ticketHandler.UpdateTicket)
		protected.DELETE("/tickets/:id", ticketHandler.DeleteTicket)
		protected.GET("/tickets/:id", ticketHandler.GetTicket)
		protected.GET("/tickets", ticketHandler.ListTicketsForCompany)

		protected.POST("/activities", activityHandler.CreateActivity)
		protected.PUT("/activities/:id", activityHandler.UpdateActivity)
		protected.DELETE("/activities/:id", activityHandler.DeleteActivity)
		protected.GET("/activities/:id", activityHandler.GetActivity)
		protected.GET("/activities", activityHandler.ListActivitiesForCompany)

		protected.POST("/entries", entryHandler.CreateEntry)
		protected.PUT("/entries/:id", entryHandler.UpdateEntry)
		protected.DELETE("/entries/:id", entryHandler.DeleteEntry)
		protected.GET("/entries/:id", entryHandler.GetEntry)

		protected.GET("/timesheet", timesheetHandler.GetTimesheet)

		protected.POST("/invoices/generate", invoiceHandler.GenerateInvoice)
	}

	return r
}
