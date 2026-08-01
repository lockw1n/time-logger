package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// parseEntryID validates the positional entry-id argument shared by edit and
// delete: a positive integer.
func parseEntryID(s string) (uint64, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid entry id %q", s)
	}
	return id, nil
}

// notFoundErr reports whether err is a 404 from the API — used to turn a raw
// "not found" into the friendly `entry N not found` message.
func notFoundErr(err error) bool {
	var apiErr *api.APIError
	return errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound
}

// entryTicket resolves an entry's ticket code via a single-day timesheet lookup
// so edit/delete prompts can name the ticket that GET /entries/:id identifies
// only by id. Best-effort: on any lookup failure it falls back to the ticket id.
func entryTicket(ctx context.Context, client api.Client, e api.Entry) string {
	fallback := fmt.Sprintf("ticket #%d", e.TicketID)
	ts, err := client.Timesheet(ctx, e.CompanyID, e.Date, e.Date)
	if err != nil {
		return fallback
	}
	for _, row := range ts.Rows {
		for _, en := range row.Entries {
			if en.ID == e.ID {
				return row.Ticket.Code
			}
		}
	}
	return fallback
}

// quoteComment renders a comment for a prompt: quoted, or "(none)" when empty.
func quoteComment(s string) string {
	if s == "" {
		return "(none)"
	}
	return strconv.Quote(s)
}
