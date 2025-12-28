package app

import (
	"os"

	"github.com/lockw1n/time-logger/internal/pdf"
	"gorm.io/gorm"

	activityrepo "github.com/lockw1n/time-logger/internal/activity/repository"
	activityservice "github.com/lockw1n/time-logger/internal/activity/service"
	companyrepo "github.com/lockw1n/time-logger/internal/company/repository"
	companyservice "github.com/lockw1n/time-logger/internal/company/service"
	consultantrepo "github.com/lockw1n/time-logger/internal/consultant/repository"
	consultantservice "github.com/lockw1n/time-logger/internal/consultant/service"
	contractrepo "github.com/lockw1n/time-logger/internal/contract/repository"
	contractservice "github.com/lockw1n/time-logger/internal/contract/service"
	entryrepo "github.com/lockw1n/time-logger/internal/entry/repository"
	entryservice "github.com/lockw1n/time-logger/internal/entry/service"
	invoiceservice "github.com/lockw1n/time-logger/internal/invoice/service"
	ticketrepo "github.com/lockw1n/time-logger/internal/ticket/repository"
	ticketservice "github.com/lockw1n/time-logger/internal/ticket/service"
	timesheetservice "github.com/lockw1n/time-logger/internal/timesheet/service"
)

type Container struct {
	DB *gorm.DB

	CompanyService    companyservice.Service
	ConsultantService consultantservice.Service
	ContractService   contractservice.Service
	TicketService     ticketservice.Service
	ActivityService   activityservice.Service
	EntryService      entryservice.Service
	TimesheetService  timesheetservice.Service
	InvoiceService    invoiceservice.Service
	PdfRenderer       pdf.Renderer
}

func NewContainer(db *gorm.DB) *Container {
	companyRepo := companyrepo.NewGormRepository(db)
	consultantRepo := consultantrepo.NewGormRepository(db)
	contractRepo := contractrepo.NewGormRepository(db)
	ticketRepo := ticketrepo.NewGormRepository(db)
	activityRepo := activityrepo.NewGormRepository(db)
	entryRepo := entryrepo.NewGormRepository(db)

	return &Container{
		DB:                db,
		CompanyService:    companyservice.NewService(companyRepo),
		ConsultantService: consultantservice.NewService(consultantRepo),
		ContractService:   contractservice.NewService(contractRepo),
		TicketService:     ticketservice.NewService(ticketRepo),
		ActivityService:   activityservice.NewService(activityRepo),
		EntryService:      entryservice.NewService(entryRepo, contractRepo, ticketRepo),
		TimesheetService:  timesheetservice.NewService(activityRepo, entryRepo, ticketRepo),
		InvoiceService:    invoiceservice.NewService(consultantRepo, companyRepo, contractRepo, activityRepo, entryRepo, ticketRepo),
		PdfRenderer:       pdf.NewHTTPRenderer(os.Getenv("PDF_RENDERER_URL"), os.Getenv("PDF_RENDERER_TOKEN")),
	}
}
