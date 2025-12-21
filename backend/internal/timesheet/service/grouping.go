package service

import (
	"sort"

	entrymapper "github.com/lockw1n/time-logger/internal/mapper/entry"
	labelmapper "github.com/lockw1n/time-logger/internal/mapper/label"
	ticketmapper "github.com/lockw1n/time-logger/internal/mapper/ticket"
	"github.com/lockw1n/time-logger/internal/models"
	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type groupKey struct {
	ticketID uint64
	labelID  uint64
}

func groupEntries(entries []models.Entry) []domain.ReportRow {
	groups := map[groupKey]*domain.ReportRow{}

	for _, e := range entries {
		var ticketID uint64
		if e.Ticket != nil {
			ticketID = e.Ticket.ID
		}

		var labelID uint64
		if e.Label != nil {
			labelID = e.Label.ID
		}

		key := groupKey{ticketID: ticketID, labelID: labelID}

		row, exists := groups[key]
		if !exists {
			row = &domain.ReportRow{
				Ticket: ticketmapper.ToResponse(e.Ticket),
				Label:  labelmapper.ToResponse(e.Label),
			}
			groups[key] = row
		}

		row.Entries = append(row.Entries, entrymapper.ToShortResponse(&e))
		row.Total += e.DurationMinutes
	}

	rows := make([]domain.ReportRow, 0, len(groups))
	for _, row := range groups {
		rows = append(rows, *row)
	}

	return rows
}

func sortRowsByTicketCode(rows []domain.ReportRow) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ticket.Code < rows[j].Ticket.Code
	})
}
