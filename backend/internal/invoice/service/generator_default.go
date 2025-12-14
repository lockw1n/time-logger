package service

import (
	"context"
	"math"
	"time"

	"github.com/lockw1n/time-logger/internal/invoice/domain"
	"github.com/lockw1n/time-logger/internal/models"
	companyrepo "github.com/lockw1n/time-logger/internal/repository/company"
	consultantrepo "github.com/lockw1n/time-logger/internal/repository/consultant"
	assignmentrepo "github.com/lockw1n/time-logger/internal/repository/consultantassignment"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
)

type invoiceGenerator struct {
	assignmentRepo assignmentrepo.Repository
	companyRepo    companyrepo.Repository
	consultantRepo consultantrepo.Repository
	entryRepo      entryrepo.Repository
	clock          Clock
}

func NewInvoiceGenerator(
	assignmentRepo assignmentrepo.Repository,
	companyRepo companyrepo.Repository,
	consultantRepo consultantrepo.Repository,
	entryRepo entryrepo.Repository,
	clock Clock,
) InvoiceGenerator {
	return &invoiceGenerator{
		assignmentRepo: assignmentRepo,
		companyRepo:    companyRepo,
		consultantRepo: consultantRepo,
		entryRepo:      entryRepo,
		clock:          clock,
	}
}

func (g *invoiceGenerator) GenerateMonthly(
	ctx context.Context,
	cmd GenerateMonthlyInvoiceCommand,
) (*domain.Invoice, error) {

	// 1. Parse month → period
	period, err := buildPeriod(cmd.Month)
	if err != nil {
		return nil, ErrInvalidMonthFormat
	}

	// 2. Resolve assignment (rate + currency source)
	assignments, err := g.assignmentRepo.FindActiveForPeriod(
		ctx,
		cmd.ConsultantID,
		cmd.CompanyID,
		period.Start,
		period.End,
	)
	if err != nil {
		return nil, err
	}

	if len(assignments) == 0 {
		return nil, ErrAssignmentNotFound
	}
	if len(assignments) > 1 {
		return nil, ErrMultipleAssignments
	}

	assignment := assignments[0]

	// 3. Load entries
	entries, err := g.entryRepo.FindForPeriodWithDetails(
		ctx,
		cmd.ConsultantID,
		cmd.CompanyID,
		period.Start,
		period.End,
	)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, ErrNoEntriesForPeriod
	}

	// 4. Load consultant + company snapshots
	consultant, err := g.consultantRepo.FindByID(ctx, cmd.ConsultantID)
	if err != nil {
		return nil, err
	}

	company, err := g.companyRepo.FindByID(ctx, cmd.CompanyID)
	if err != nil {
		return nil, err
	}

	// 5. Build groups + totals
	groups := buildGroups(entries, assignment)
	totals := calculateTotals(groups)

	// 6. Assemble domain invoice
	now := g.clock.Now()

	invoice := &domain.Invoice{
		Number:   generateInvoiceNumber(cmd.Month),
		IssuedAt: now,
		DueAt:    now.AddDate(0, 0, 14),
		Currency: assignment.Currency,

		Period:     period,
		Consultant: mapConsultant(*consultant),
		Company:    mapCompany(*company),
		Contract:   mapContract(assignment),

		Groups: groups,
		Totals: totals,
	}

	return invoice, nil
}

func buildPeriod(month string) (domain.Period, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return domain.Period{}, ErrInvalidMonthFormat
	}

	start := time.Date(
		t.Year(),
		t.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	// First day of next month at midnight
	end := start.AddDate(0, 1, 0)

	return domain.Period{
		Month: month,
		Start: start,
		End:   end,
	}, nil
}

func generateInvoiceNumber(month string) string {
	// Parse YYYY-MM
	t, err := time.Parse("2006-01", month)
	if err != nil {
		// This should never happen if buildPeriod already validated,
		panic("invalid month format: " + month)
	}

	// First day of the next month
	nextMonthFirstDay := time.Date(
		t.Year(),
		t.Month()+1,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	return nextMonthFirstDay.Format("20060102")
}

func buildGroups(
	entries []models.Entry,
	assignment models.ConsultantAssignment,
) []domain.Group {

	groupsMap := make(map[string]*domain.Group)

	for _, entry := range entries {
		label := "Unlabeled"
		if entry.Label != nil && entry.Label.Name != "" {
			label = entry.Label.Name
		}

		// Minutes → hours
		hours := float64(entry.DurationMinutes) / 60.0

		// Resolve hourly rate (snapshot preferred)
		hourlyRate := assignment.HourlyRate
		if entry.HourlyRateSnapshot != nil {
			hourlyRate = *entry.HourlyRateSnapshot
		}

		// Amount in smallest unit (assume 2 decimals)
		amount := int64(math.Round(hours * hourlyRate * 100))

		row := domain.Row{
			Date:        entry.Date,
			TicketCode:  ticketCode(entry),
			Description: entryDescription(entry),
			Hours:       hours,
			Amount:      amount,
		}

		group, exists := groupsMap[label]
		if !exists {
			group = &domain.Group{
				Label: label,
			}
			groupsMap[label] = group
		}

		group.Rows = append(group.Rows, row)
		group.HourlyRate = hourlyRate
		group.TotalHours += hours
		group.Subtotal += amount
	}

	groups := make([]domain.Group, 0, len(groupsMap))
	for _, g := range groupsMap {
		groups = append(groups, *g)
	}

	return groups
}

func ticketCode(entry models.Entry) string {
	if entry.Ticket != nil {
		return entry.Ticket.Code
	}
	return ""
}

func entryDescription(entry models.Entry) string {
	if entry.Comment != nil {
		return *entry.Comment
	}
	if entry.Ticket != nil && entry.Ticket.Description != nil {
		return *entry.Ticket.Description
	}
	return ""
}

func calculateTotals(groups []domain.Group) domain.Totals {
	var totalHours float64
	var subtotal int64

	for _, group := range groups {
		totalHours += group.TotalHours
		subtotal += group.Subtotal
	}

	return domain.Totals{
		TotalHours: totalHours,
		Subtotal:   subtotal,
	}
}

func mapConsultant(consultant models.Consultant) domain.Consultant {
	return domain.Consultant{
		FirstName:  consultant.FirstName,
		MiddleName: derefString(consultant.MiddleName),
		LastName:   consultant.LastName,

		AddressLine1: consultant.AddressLine1,
		AddressLine2: derefString(consultant.AddressLine2),
		Zip:          consultant.Zip,
		City:         consultant.City,
		Region:       derefString(consultant.Region),
		Country:      consultant.Country,

		TaxNumber: derefString(consultant.TaxNumber),

		BankName:    consultant.BankName,
		BankAddress: consultant.BankAddress,
		BankCountry: consultant.BankCountry,
		BankIban:    consultant.BankIBAN,
		BankBic:     consultant.BankBIC,
	}
}

func mapCompany(company models.Company) domain.Company {
	return domain.Company{
		Name:      company.Name,
		NameShort: derefString(company.NameShort),
		TaxNumber: derefString(company.TaxNumber),

		AddressLine1: company.AddressLine1,
		AddressLine2: derefString(company.AddressLine2),
		Zip:          company.Zip,
		City:         company.City,
		Region:       derefString(company.Region),
		Country:      company.Country,

		PaymentTerms: derefString(company.PaymentTerms),
	}
}

func mapContract(assignment models.ConsultantAssignment) domain.Contract {
	return domain.Contract{
		OrderNumber: assignment.OrderNumber,
		HourlyRate:  assignment.HourlyRate,
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
