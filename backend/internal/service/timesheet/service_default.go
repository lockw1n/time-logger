package timesheet

import (
	"fmt"
	"sort"
	"time"

	"github.com/lockw1n/time-logger/internal/constants"
	entrydto "github.com/lockw1n/time-logger/internal/dto/entry"
	timesheetdto "github.com/lockw1n/time-logger/internal/dto/timesheet"
	entrymapper "github.com/lockw1n/time-logger/internal/mapper/entry"
	labelmapper "github.com/lockw1n/time-logger/internal/mapper/label"
	ticketmapper "github.com/lockw1n/time-logger/internal/mapper/ticket"
	entryrepo "github.com/lockw1n/time-logger/internal/repository/entry"
)

type service struct {
	entryRepo entryrepo.Repository
}

func NewService(entryRepo entryrepo.Repository) Service {
	return &service{entryRepo: entryRepo}
}

func makeGroupKey(ticketID, labelID uint64) string {
	return fmt.Sprintf("%d-%d", ticketID, labelID)
}

func (s *service) GenerateReport(req timesheetdto.Request) (*timesheetdto.Report, error) {
	// Parse dates
	start, err := time.Parse(constants.DateFormat, req.Start)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(constants.DateFormat, req.End)
	if err != nil {
		return nil, err
	}

	// Load enriched entries: includes Ticket + Label
	entries, err := s.entryRepo.FindWithDetails(req.ConsultantID, req.CompanyID, start, end)
	if err != nil {
		return nil, err
	}

	// Group by (ticket_id, label_id)
	groups := map[string]*timesheetdto.ReportRow{}

	for _, e := range entries {
		var ticketID uint64
		var labelID uint64

		if e.Ticket != nil {
			ticketID = e.Ticket.ID
		} else {
			ticketID = 0
		}

		if e.Label != nil {
			labelID = e.Label.ID
		} else {
			labelID = 0
		}

		key := makeGroupKey(ticketID, labelID)

		// Initialize group if not exist
		if _, exists := groups[key]; !exists {
			groups[key] = &timesheetdto.ReportRow{
				Ticket:  ticketmapper.ToResponse(e.Ticket),
				Label:   labelmapper.ToResponse(e.Label),
				Entries: []entrydto.ShortResponse{},
				Total:   0,
			}
		}

		// Append entry to the group
		groups[key].Entries = append(groups[key].Entries, entrymapper.ToShortResponse(&e))

		groups[key].Total += e.DurationMinutes
	}

	// Convert map → slice
	rows := make([]timesheetdto.ReportRow, 0, len(groups))
	for _, row := range groups {
		rows = append(rows, *row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ticket.Code < rows[j].Ticket.Code
	})

	// Build final report
	return &timesheetdto.Report{
		ConsultantID: req.ConsultantID,
		CompanyID:    req.CompanyID,
		Start:        req.Start,
		End:          req.End,
		Rows:         rows,
	}, nil
}
