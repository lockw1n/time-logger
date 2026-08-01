package api

import "time"

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Deref returns the pointed-to string, or "" for a nil pointer. It lives here
// because the response DTOs use *string for nullable fields, and both the
// render and command layers need to flatten them for display.
func Deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Company mirrors the backend CompanyResponse. All fields are declared so that
// `tl companies --json` round-trips the payload faithfully.
type Company struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	NameShort    *string `json:"name_short"`
	TaxNumber    string  `json:"tax_number"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Zip          string  `json:"zip"`
	City         string  `json:"city"`
	Region       *string `json:"region"`
	Country      string  `json:"country"`
}

// Activity mirrors the backend ActivityResponse.
type Activity struct {
	ID        uint64  `json:"id"`
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	Billable  bool    `json:"billable"`
	Priority  int     `json:"priority"`
}

// CreateEntryRequest mirrors the backend CreateEntryRequest. The server derives
// the consultant from the JWT, so consultant_id is intentionally omitted.
type CreateEntryRequest struct {
	CompanyID       uint64  `json:"company_id"`
	TicketCode      string  `json:"ticket_code"`
	ActivityID      uint64  `json:"activity_id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

// UpdateEntryRequest mirrors the backend UpdateEntryRequest — only duration and
// comment are updatable.
type UpdateEntryRequest struct {
	DurationMinutes *int    `json:"duration_minutes,omitempty"`
	Comment         *string `json:"comment,omitempty"`
}

// Entry mirrors the backend EntryResponse.
type Entry struct {
	ID              uint64  `json:"id"`
	ConsultantID    uint64  `json:"consultant_id"`
	CompanyID       uint64  `json:"company_id"`
	ContractID      uint64  `json:"contract_id"`
	TicketID        uint64  `json:"ticket_id"`
	ActivityID      uint64  `json:"activity_id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

// GenerateInvoiceRequest carries the query parameters for POST
// /api/invoices/generate. Format is "pdf" or "excel"; dates are YYYY-MM-DD.
type GenerateInvoiceRequest struct {
	CompanyID uint64
	Start     string
	End       string
	Format    string
}

// InvoiceFile is a generated invoice document: the raw bytes plus the filename
// the client resolved (from Content-Disposition, or a start/end fallback).
type InvoiceFile struct {
	Filename string
	Data     []byte
}

// TimesheetResponse mirrors the backend timesheet payload.
type TimesheetResponse struct {
	ConsultantID uint64          `json:"consultant_id"`
	CompanyID    uint64          `json:"company_id"`
	Start        string          `json:"start"`
	End          string          `json:"end"`
	Rows         []TimesheetRow  `json:"rows"`
	Totals       TimesheetTotals `json:"totals"`
}

type TimesheetRow struct {
	Ticket        TimesheetTicket   `json:"ticket"`
	Activity      TimesheetActivity `json:"activity"`
	Entries       []TimesheetEntry  `json:"entries"`
	PerDayMinutes map[string]int    `json:"per_day_minutes"`
	TotalMinutes  int               `json:"total_minutes"`
}

type TimesheetTicket struct {
	ID          uint64  `json:"id"`
	CompanyID   uint64  `json:"company_id"`
	Code        string  `json:"code"`
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

type TimesheetActivity struct {
	ID        uint64  `json:"id"`
	CompanyID uint64  `json:"company_id"`
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	Billable  bool    `json:"billable"`
	Priority  int     `json:"priority"`
}

type TimesheetEntry struct {
	ID              uint64  `json:"id"`
	Date            string  `json:"date"`
	DurationMinutes int     `json:"duration_minutes"`
	Comment         *string `json:"comment"`
}

type TimesheetTotals struct {
	PerDayMinutes  map[string]int `json:"per_day_minutes"`
	OverallMinutes int            `json:"overall_minutes"`
}
