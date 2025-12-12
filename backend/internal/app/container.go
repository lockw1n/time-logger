package app

import (
	"gorm.io/gorm"

	companyrepo "github.com/lockw1n/time-logger/internal/repository/company"
	consultantrepo "github.com/lockw1n/time-logger/internal/repository/consultant"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
	invoicerepo "github.com/lockw1n/time-logger/internal/repository/invoice"
	invoicelinerepo "github.com/lockw1n/time-logger/internal/repository/invoiceline"
	labelrepo "github.com/lockw1n/time-logger/internal/repository/label"
	ticketrepo "github.com/lockw1n/time-logger/internal/repository/ticket"
	companyservice "github.com/lockw1n/time-logger/internal/service/company"
	consultantservice "github.com/lockw1n/time-logger/internal/service/consultant"
	assignmentservice "github.com/lockw1n/time-logger/internal/service/consultantassignment"
	entryservice "github.com/lockw1n/time-logger/internal/service/entry"
	invoiceservice "github.com/lockw1n/time-logger/internal/service/invoice"
	labelservice "github.com/lockw1n/time-logger/internal/service/label"
	ticketservice "github.com/lockw1n/time-logger/internal/service/ticket"
	timesheetService "github.com/lockw1n/time-logger/internal/service/timesheet"
)

type Container struct {
	DB *gorm.DB

	CompanyService              companyservice.Service
	ConsultantService           consultantservice.Service
	ConsultantAssignmentService assignmentservice.Service
	EntryService                entryservice.Service
	LabelService                labelservice.Service
	TicketService               ticketservice.Service
	InvoiceService              invoiceservice.Service
	TimesheetService            timesheetService.Service
}

func NewContainer(db *gorm.DB) *Container {
	companyRepo := companyrepo.NewGormRepository(db)
	consultantRepo := consultantrepo.NewGormRepository(db)
	assignmentRepo := assignmentrepo.NewGormRepository(db)
	entryRepo := entryrepo.NewGormRepository(db)
	labelRepo := labelrepo.NewGormRepository(db)
	ticketRepo := ticketrepo.NewGormRepository(db)
	invoiceRepo := invoicerepo.NewGormRepository(db)
	invoiceLineRepo := invoicelinerepo.NewGormRepository(db)

	return &Container{
		DB:                          db,
		CompanyService:              companyservice.NewService(companyRepo),
		ConsultantService:           consultantservice.NewService(consultantRepo),
		ConsultantAssignmentService: assignmentservice.NewService(assignmentRepo),
		EntryService:                entryservice.NewService(entryRepo, assignmentRepo, ticketRepo),
		LabelService:                labelservice.NewService(labelRepo),
		TicketService:               ticketservice.NewService(ticketRepo),
		InvoiceService:              invoiceservice.NewService(entryRepo, invoiceRepo, invoiceLineRepo, assignmentRepo),
		TimesheetService:            timesheetService.NewService(entryRepo),
	}
}
