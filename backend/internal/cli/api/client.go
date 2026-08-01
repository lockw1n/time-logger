package api

import (
	"context"
	"encoding/json"
)

// Client is the seam every command depends on to reach the backend. The
// concrete implementation (httpClient) talks HTTP+JSON through nginx; tests and
// later phases (e.g. the offline outbox) substitute their own implementations.
type Client interface {
	Login(ctx context.Context, email, password string) (LoginResult, error)
	Health(ctx context.Context) error

	Companies(ctx context.Context) ([]Company, error)
	Activities(ctx context.Context, companyID uint64) ([]Activity, error)
	Timesheet(ctx context.Context, companyID uint64, start, end string) (TimesheetResponse, error)
	CreateEntry(ctx context.Context, req CreateEntryRequest) (Entry, error)
	UpdateEntry(ctx context.Context, id uint64, req UpdateEntryRequest) (Entry, error)
	GetEntry(ctx context.Context, id uint64) (Entry, error)
	DeleteEntry(ctx context.Context, id uint64) error

	// GenerateInvoice streams an invoice document (PDF or Excel) for a period.
	GenerateInvoice(ctx context.Context, req GenerateInvoiceRequest) (InvoiceFile, error)

	// GetRaw returns an endpoint's response body verbatim — used by --json.
	GetRaw(ctx context.Context, path string) (json.RawMessage, error)
}
