package app

import (
	"gorm.io/gorm"

	"github.com/lockw1n/time-logger/internal/repository/company"
	"github.com/lockw1n/time-logger/internal/repository/consultant"
	"github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	"github.com/lockw1n/time-logger/internal/repository/entry"
	"github.com/lockw1n/time-logger/internal/repository/invoice"
	"github.com/lockw1n/time-logger/internal/repository/invoiceline"
	"github.com/lockw1n/time-logger/internal/repository/label"
	"github.com/lockw1n/time-logger/internal/repository/ticket"

	svcCompany "github.com/lockw1n/time-logger/internal/service/company"
	svcConsultant "github.com/lockw1n/time-logger/internal/service/consultant"
	svcAssignment "github.com/lockw1n/time-logger/internal/service/consultantassignment"
	svcEntry "github.com/lockw1n/time-logger/internal/service/entry"
	svcInvoice "github.com/lockw1n/time-logger/internal/service/invoice"
	svcLabel "github.com/lockw1n/time-logger/internal/service/label"
	svcTicket "github.com/lockw1n/time-logger/internal/service/ticket"
	svcTimesheet "github.com/lockw1n/time-logger/internal/service/timesheet"
)

type Container struct {
	DB *gorm.DB

	CompanyService              svcCompany.Service
	ConsultantService           svcConsultant.Service
	ConsultantAssignmentService svcAssignment.Service
	EntryService                svcEntry.Service
	LabelService                svcLabel.Service
	TicketService               svcTicket.Service
	InvoiceService              svcInvoice.Service
	TimesheetService            svcTimesheet.Service
}

func NewContainer(db *gorm.DB) *Container {
	companyRepo := company.NewGormRepository(db)
	consultantRepo := consultant.NewGormRepository(db)
	assignmentRepo := consultantassignment.NewGormRepository(db)
	entryRepo := entry.NewGormRepository(db)
	labelRepo := label.NewGormRepository(db)
	ticketRepo := ticket.NewGormRepository(db)
	invoiceRepo := invoice.NewGormRepository(db)
	invoiceLineRepo := invoiceline.NewGormRepository(db)

	return &Container{
		DB:                          db,
		CompanyService:              svcCompany.NewService(companyRepo),
		ConsultantService:           svcConsultant.NewService(consultantRepo),
		ConsultantAssignmentService: svcAssignment.NewService(assignmentRepo),
		EntryService:                svcEntry.NewService(entryRepo, assignmentRepo, ticketRepo),
		LabelService:                svcLabel.NewService(labelRepo),
		TicketService:               svcTicket.NewService(ticketRepo),
		InvoiceService:              svcInvoice.NewService(entryRepo, invoiceRepo, invoiceLineRepo, assignmentRepo),
		TimesheetService:            svcTimesheet.NewService(entryRepo),
	}
}
