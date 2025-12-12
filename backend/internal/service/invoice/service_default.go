package invoice

import (
	"errors"
	"time"

	invoicedto "github.com/lockw1n/time-logger/internal/dto/invoice"
	invoicelinedto "github.com/lockw1n/time-logger/internal/dto/invoiceline"
	invoicemapper "github.com/lockw1n/time-logger/internal/mapper/invoice"
	invoicelinemapper "github.com/lockw1n/time-logger/internal/mapper/invoiceline"
	"github.com/lockw1n/time-logger/internal/models"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
	invoicerepo "github.com/lockw1n/time-logger/internal/repository/invoice"
	invoicelinerepo "github.com/lockw1n/time-logger/internal/repository/invoiceline"
)

type service struct {
	entryRepo       entryrepo.Repository
	invoiceRepo     invoicerepo.Repository
	invoiceLineRepo invoicelinerepo.Repository
	assignmentRepo  assignmentrepo.Repository
}

func NewService(
	entryRepo entryrepo.Repository,
	invoiceRepo invoicerepo.Repository,
	invoiceLineRepo invoicelinerepo.Repository,
	assignmentRepo assignmentrepo.Repository,
) Service {
	return &service{
		entryRepo:       entryRepo,
		invoiceRepo:     invoiceRepo,
		invoiceLineRepo: invoiceLineRepo,
		assignmentRepo:  assignmentRepo,
	}
}

/* ============================================================
   DYNAMIC INVOICE PREVIEW (NOT SAVED IN DB)
=============================================================== */

func (s *service) Generate(consultantID, companyID uint64, start, end string) (*invoicedto.GenerateResponse, error) {
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	// Fetch consultant's entries
	entries, err := s.entryRepo.FindByConsultant(consultantID)
	if err != nil {
		return nil, err
	}

	// Filter for company + date range
	filtered := make([]models.Entry, 0)
	for _, e := range entries {
		if e.CompanyID == companyID &&
			!e.Date.Before(startDate) &&
			!e.Date.After(endDate) {
			filtered = append(filtered, e)
		}
	}

	// Get hourly rate
	cc, err := s.assignmentRepo.FindByPair(consultantID, companyID)
	if err != nil {
		return nil, errors.New("hourly rate not found for consultant/company pair")
	}

	// Convert entries → invoice line DTOs
	lineDTOs := make([]invoicelinedto.Line, len(filtered))
	var total float64

	for i := range filtered {
		e := filtered[i]

		hours := float64(e.DurationMinutes) / 60.0
		amount := hours * cc.HourlyRate

		lineDTOs[i] = invoicelinedto.Line{
			ID:        0, // preview only — not stored
			EntryID:   e.ID,
			Hours:     hours,
			Rate:      cc.HourlyRate,
			Amount:    amount,
			CreatedAt: time.Now().Format(time.RFC3339), // preview only
		}

		total += amount
	}

	// Build response via mapper
	out := invoicemapper.ToGenerateResponse(
		consultantID,
		companyID,
		start,
		end,
		total,
		cc.Currency,
		lineDTOs,
	)

	return &out, nil
}

/* ============================================================
   LIST INVOICES FOR COMPANY
=============================================================== */

func (s *service) ListByCompany(companyID uint64) ([]invoicedto.ListItem, error) {
	list, err := s.invoiceRepo.FindByCompany(companyID)
	if err != nil {
		return nil, err
	}

	return invoicemapper.ToListItems(list), nil
}

/* ============================================================
   LIST INVOICES FOR CONSULTANT
=============================================================== */

func (s *service) ListByConsultant(consultantID uint64) ([]invoicedto.ListItem, error) {
	list, err := s.invoiceRepo.FindByConsultant(consultantID)
	if err != nil {
		return nil, err
	}

	return invoicemapper.ToListItems(list), nil
}

/* ============================================================
   GET INVOICE DETAILS
=============================================================== */

func (s *service) Get(id uint64) (*invoicedto.Response, error) {
	inv, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Load invoice lines
	lines, err := s.invoiceLineRepo.FindByInvoice(id)
	if err != nil {
		return nil, err
	}

	dtoLines := invoicelinemapper.ToLines(lines)

	out := invoicemapper.ToResponse(inv)
	out.Lines = dtoLines

	return &out, nil
}
