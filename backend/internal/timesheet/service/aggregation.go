package service

import (
	"fmt"
	"sort"

	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	entrydomain "github.com/lockw1n/time-logger/internal/entry/domain"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
	"github.com/lockw1n/time-logger/internal/timesheet/domain"
)

type groupKey struct {
	TicketID   uint64
	ActivityID uint64
}

func groupEntries(
	entries []entrydomain.Entry,
	ticketsByID map[uint64]ticketdomain.Ticket,
	activitiesByID map[uint64]activitydomain.Activity,
) []domain.TimesheetRow {

	rowsByKey := make(map[groupKey]*domain.TimesheetRow)

	for _, entry := range entries {
		key := groupKey{
			TicketID:   entry.TicketID,
			ActivityID: entry.ActivityID,
		}

		row, exists := rowsByKey[key]
		if !exists {
			ticket, ok := ticketsByID[entry.TicketID]
			if !ok {
				panic(fmt.Sprintf(
					"timesheet integrity violation: entry %d references missing ticket %d",
					entry.ID,
					entry.TicketID,
				))
			}

			activity, ok := activitiesByID[entry.ActivityID]
			if !ok {
				panic(fmt.Sprintf(
					"timesheet integrity violation: entry %d references missing activity %d",
					entry.ID,
					entry.ActivityID,
				))
			}

			row = &domain.TimesheetRow{
				Ticket:   ticket,
				Activity: activity,
				Entries:  make([]domain.TimesheetEntry, 0),
				Total:    0,
			}

			rowsByKey[key] = row
		}

		row.Entries = append(row.Entries, domain.TimesheetEntry{
			ID:              entry.ID,
			Date:            entry.Date,
			DurationMinutes: entry.DurationMinutes,
			Comment:         entry.Comment,
		})

		row.Total += entry.DurationMinutes
	}

	rows := make([]domain.TimesheetRow, 0, len(rowsByKey))
	for _, row := range rowsByKey {
		rows = append(rows, *row)
	}

	return rows
}

func sortRowsByTicketCode(rows []domain.TimesheetRow) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ticket.Code < rows[j].Ticket.Code
	})
}
