// Package entrylog implements the shared "submit a time entry" flow used by both
// `tl add` and `tl stop`: POST the entry, and when the backend reports a
// duplicate (409) for the same ticket/activity/day, offer to merge by summing
// the minutes into the existing entry. Keeping it here means the timer's stop
// path and the manual add path behave identically instead of drifting apart.
package entrylog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

// ErrCreateUnreachable reports that the *initial* create never reached the
// backend, so the caller may safely queue the entry for offline replay. It is
// returned ONLY from the create attempt — never from the merge path, where an
// UpdateEntry may already have committed and re-queuing a fresh create would
// double-count. It wraps api.ErrUnreachable so the outbox replay classifier still
// treats it as a retryable transient failure.
var ErrCreateUnreachable = fmt.Errorf("%w: create did not reach the backend", api.ErrUnreachable)

// Request describes one entry to submit. The activity is already resolved so the
// caller (which knows how to look up and validate names) owns that step.
type Request struct {
	CompanyID  uint64
	TicketCode string
	Activity   api.Activity
	Date       string
	Minutes    int
	Comment    *string
}

// Options carries the caller's I/O: where to print the human summary and how to
// ask the merge confirmation. Confirm is skipped entirely when SkipConfirm is set.
type Options struct {
	Out         io.Writer
	Confirm     func(prompt string, defaultYes bool) bool
	SkipConfirm bool
}

// Submit creates the entry, transparently handling the duplicate-merge case. On
// success it writes a one-line summary to opts.Out. A merge declined at the
// prompt returns an "aborted" error and changes nothing.
func Submit(ctx context.Context, client api.Client, req Request, opts Options) error {
	entry, err := client.CreateEntry(ctx, api.CreateEntryRequest{
		CompanyID:       req.CompanyID,
		TicketCode:      req.TicketCode,
		ActivityID:      req.Activity.ID,
		Date:            req.Date,
		DurationMinutes: req.Minutes,
		Comment:         req.Comment,
	})

	var apiErr *api.APIError
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusConflict {
		return merge(ctx, client, req, opts)
	}
	if err != nil {
		// A create that never landed is safe to queue offline; tag it so the
		// caller can tell it apart from a merge-path failure (which must not be
		// re-queued). Merge-path unreachable errors deliberately do NOT carry
		// this sentinel.
		if errors.Is(err, api.ErrUnreachable) {
			return fmt.Errorf("%w (%v)", ErrCreateUnreachable, err)
		}
		return err
	}

	fmt.Fprintf(opts.Out, "logged %s on %s (%s) for %s — entry #%d\n",
		parse.FormatMinutes(req.Minutes), req.TicketCode, req.Activity.Name, req.Date, entry.ID)
	return nil
}

// merge disambiguates the two meanings of a 409 on POST /entries: a duplicate
// (ticket, activity, date) — in which case it offers to sum the minutes into the
// existing entry — versus no active contract for the date.
func merge(ctx context.Context, client api.Client, req Request, opts Options) error {
	ts, err := client.Timesheet(ctx, req.CompanyID, req.Date, req.Date)
	if err != nil {
		return err
	}

	existing, ok := findDuplicate(ts, req.TicketCode, req.Activity.ID, req.Date)
	if !ok {
		// The 409 was not a duplicate, so it means the date is uncovered.
		return fmt.Errorf("no active contract covers %s for company %d", req.Date, req.CompanyID)
	}

	newTotal := existing.DurationMinutes + req.Minutes
	if newTotal > parse.MaxMinutes {
		return fmt.Errorf("merged total %s exceeds the 24h per-entry limit (%s already logged + %s)",
			parse.FormatMinutes(newTotal), parse.FormatMinutes(existing.DurationMinutes), parse.FormatMinutes(req.Minutes))
	}

	if !opts.SkipConfirm {
		prompt := fmt.Sprintf("%s/%s already has %s on %s — add %s for a total of %s?",
			req.TicketCode, req.Activity.Name,
			parse.FormatMinutes(existing.DurationMinutes), req.Date,
			parse.FormatMinutes(req.Minutes), parse.FormatMinutes(newTotal))
		// A nil Confirm means the caller offered no way to answer, so a merge that
		// needs confirmation cannot proceed — abort rather than nil-panic.
		if opts.Confirm == nil || !opts.Confirm(prompt, true) {
			return errors.New("aborted")
		}
	}

	merged := mergeComments(existing.Comment, req.Comment)
	updated, err := client.UpdateEntry(ctx, existing.ID, api.UpdateEntryRequest{
		DurationMinutes: &newTotal,
		Comment:         merged,
	})
	if err != nil {
		// The merge-update may or may not have committed if the backend became
		// unreachable or the request was interrupted mid-flight. Keep the error
		// retryable for the replay classifier (it still wraps ErrUnreachable via
		// %w) but do NOT let it be re-queued as a fresh create — that would
		// double-count. The added guidance points the user at a manual check.
		if errors.Is(err, api.ErrUnreachable) || errors.Is(err, api.ErrInterrupted) {
			return fmt.Errorf("%w while merging into entry #%d — it may already be updated; check 'tl today' before retrying", err, existing.ID)
		}
		return err
	}

	fmt.Fprintf(opts.Out, "logged %s on %s (%s) for %s — entry #%d now totals %s\n",
		parse.FormatMinutes(req.Minutes), req.TicketCode, req.Activity.Name, req.Date,
		updated.ID, parse.FormatMinutes(updated.DurationMinutes))
	return nil
}

// findDuplicate locates the entry matching the ticket code (case-insensitive)
// and activity id on the given date within a single-day timesheet.
func findDuplicate(ts api.TimesheetResponse, ticketCode string, activityID uint64, date string) (api.TimesheetEntry, bool) {
	for _, row := range ts.Rows {
		if !strings.EqualFold(row.Ticket.Code, ticketCode) || row.Activity.ID != activityID {
			continue
		}
		for _, e := range row.Entries {
			if e.Date == date {
				return e, true
			}
		}
	}
	return api.TimesheetEntry{}, false
}

// mergeComments combines the existing and new comments: both present → joined
// with "; "; otherwise whichever is non-empty (nil when neither is).
func mergeComments(existing, incoming *string) *string {
	e, i := api.Deref(existing), api.Deref(incoming)
	switch {
	case e != "" && i != "":
		joined := e + "; " + i
		return &joined
	case e != "":
		return &e
	case i != "":
		return &i
	default:
		return nil
	}
}
