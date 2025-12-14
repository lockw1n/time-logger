package app

import (
	"gorm.io/gorm"

	invoiceservice "github.com/lockw1n/time-logger/internal/invoice/service"
	companyrepo "github.com/lockw1n/time-logger/internal/repository/company"
	consultantrepo "github.com/lockw1n/time-logger/internal/repository/consultant"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
	labelrepo "github.com/lockw1n/time-logger/internal/repository/label"
	ticketrepo "github.com/lockw1n/time-logger/internal/repository/ticket"
	companyservice "github.com/lockw1n/time-logger/internal/service/company"
	consultantservice "github.com/lockw1n/time-logger/internal/service/consultant"
	assignmentservice "github.com/lockw1n/time-logger/internal/service/consultantassignment"
	entryservice "github.com/lockw1n/time-logger/internal/service/entry"
	labelservice "github.com/lockw1n/time-logger/internal/service/label"
	ticketservice "github.com/lockw1n/time-logger/internal/service/ticket"
	timesheetservice "github.com/lockw1n/time-logger/internal/service/timesheet"
)

type Container struct {
	DB *gorm.DB

	CompanyService              companyservice.Service
	ConsultantService           consultantservice.Service
	ConsultantAssignmentService assignmentservice.Service
	EntryService                entryservice.Service
	LabelService                labelservice.Service
	TicketService               ticketservice.Service
	TimesheetService            timesheetservice.Service
	InvoiceGenerator            invoiceservice.InvoiceGenerator
}

func NewContainer(db *gorm.DB) *Container {
	companyRepo := companyrepo.NewGormRepository(db)
	consultantRepo := consultantrepo.NewGormRepository(db)
	assignmentRepo := assignmentrepo.NewGormRepository(db)
	entryRepo := entryrepo.NewGormRepository(db)
	labelRepo := labelrepo.NewGormRepository(db)
	ticketRepo := ticketrepo.NewGormRepository(db)
	clock := invoiceservice.NewClock()

	return &Container{
		DB:                          db,
		CompanyService:              companyservice.NewService(companyRepo),
		ConsultantService:           consultantservice.NewService(consultantRepo),
		ConsultantAssignmentService: assignmentservice.NewService(assignmentRepo),
		EntryService:                entryservice.NewService(entryRepo, assignmentRepo, ticketRepo),
		LabelService:                labelservice.NewService(labelRepo),
		TicketService:               ticketservice.NewService(ticketRepo),
		TimesheetService:            timesheetservice.NewService(entryRepo),
		InvoiceGenerator:            invoiceservice.NewInvoiceGenerator(assignmentRepo, companyRepo, consultantRepo, entryRepo, clock),
	}
}
