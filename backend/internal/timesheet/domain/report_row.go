package domain

import (
	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
	labeldto "github.com/lockw1n/time-logger/internal/dto/label"
	ticketdto "github.com/lockw1n/time-logger/internal/dto/ticket"
)

type ReportRow struct {
	Ticket  ticketdto.Response       `json:"ticket"`
	Label   labeldto.Response        `json:"label"`
	Entries []entrydto.ShortResponse `json:"entries"`
	Total   int                      `json:"total"`
}
