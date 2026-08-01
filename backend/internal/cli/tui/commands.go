package tui

import (
	"context"
	"errors"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
	"github.com/lockw1n/time-logger/internal/cli/timer"
)

// The messages below are the results of the asynchronous tea.Cmds. Every network
// call runs in a Cmd and reports back as one of these, so Update stays pure and
// the UI never blocks on I/O. Each carries an err the update loop renders into
// the footer instead of crashing the program.

// timesheetMsg delivers a (re)fetched week for the request that asked for it. The
// start guards against a stale response: if the user has since paged to another
// week or switched company, a late reply for the old view is dropped.
type timesheetMsg struct {
	companyID uint64
	start     string
	ts        api.TimesheetResponse
	err       error
}

// companiesMsg delivers the company list used to cycle with `c`.
type companiesMsg struct {
	companies []api.Company
	err       error
}

// activitiesMsg delivers the activity list used by the add form, scoped to the
// company it was fetched for so a stale list from a previous company is ignored.
type activitiesMsg struct {
	companyID  uint64
	activities []api.Activity
	err        error
}

// mutationMsg reports the outcome of a write (create/merge, update, delete). verb
// names the action for the footer; on success the update loop refetches the week
// to reconcile, on error it shows err and leaves the grid as it was. queued marks
// a create that the backend couldn't reach and that was persisted to the offline
// outbox instead — a success that must NOT trigger a reload (still offline).
type mutationMsg struct {
	verb   string
	queued bool
	err    error
}

// tickMsg drives the once-per-second header refresh. On each tick the model
// re-arms the tick and fires refreshHeader to reload the cached status.
type tickMsg time.Time

// headerMsg carries a freshly-read snapshot of the timer/outbox status for the
// header. It is produced by refreshHeader so the file reads happen in a command
// (once per tick) rather than on every render.
type headerMsg headerInfo

// fetchTimesheet loads the week [start,end] for company. The response is tagged
// with companyID/start so a slow reply that arrives after the user moved on can
// be discarded by Update.
func fetchTimesheet(ctx context.Context, client api.Client, companyID uint64, start, end string) tea.Cmd {
	return func() tea.Msg {
		ts, err := client.Timesheet(ctx, companyID, start, end)
		return timesheetMsg{companyID: companyID, start: start, ts: ts, err: err}
	}
}

// fetchCompanies loads the consultant's companies for the `c` switcher.
func fetchCompanies(ctx context.Context, client api.Client) tea.Cmd {
	return func() tea.Msg {
		companies, err := client.Companies(ctx)
		return companiesMsg{companies: companies, err: err}
	}
}

// fetchActivities loads the activities for company, used to populate the add form.
func fetchActivities(ctx context.Context, client api.Client, companyID uint64) tea.Cmd {
	return func() tea.Msg {
		activities, err := client.Activities(ctx, companyID)
		return activitiesMsg{companyID: companyID, activities: activities, err: err}
	}
}

// updateDuration PUTs a new minute count for a single entry — the single-entry
// cell edit and the entry-list edit both land here.
func updateDuration(ctx context.Context, client api.Client, id uint64, minutes int) tea.Cmd {
	return func() tea.Msg {
		_, err := client.UpdateEntry(ctx, id, api.UpdateEntryRequest{DurationMinutes: &minutes})
		return mutationMsg{verb: "updated", err: err}
	}
}

// deleteEntry DELETEs a single entry.
func deleteEntry(ctx context.Context, client api.Client, id uint64) tea.Cmd {
	return func() tea.Msg {
		err := client.DeleteEntry(ctx, id)
		return mutationMsg{verb: "deleted", err: err}
	}
}

// createEntry submits a new entry through the shared entrylog seam, so adding to
// a day where the row already has an entry goes through the same create-or-merge
// path as `tl add` — the server's unique index is never violated and minutes are
// summed instead of duplicated. Confirmation is skipped: the merge is implicit,
// the UI having already put the cursor on the row. When the backend is
// unreachable it queues the create to the offline outbox — matching `tl add` —
// so an entry typed in the UI while offline is never lost.
func createEntry(ctx context.Context, client api.Client, req entrylog.Request) tea.Cmd {
	return func() tea.Msg {
		err := entrylog.Submit(ctx, client, req, entrylog.Options{Out: io.Discard, SkipConfirm: true})
		if errors.Is(err, entrylog.ErrCreateUnreachable) {
			if _, qErr := outbox.Enqueue(outbox.Op{
				Kind:      outbox.KindCreateEntry,
				CreatedAt: time.Now(),
				Payload: api.CreateEntryRequest{
					CompanyID:       req.CompanyID,
					TicketCode:      req.TicketCode,
					ActivityID:      req.Activity.ID,
					Date:            req.Date,
					DurationMinutes: req.Minutes,
					Comment:         req.Comment,
				},
				MergeOnConflict: true,
			}); qErr != nil {
				return mutationMsg{verb: "add", err: qErr}
			}
			return mutationMsg{verb: "queued for sync (offline)", queued: true}
		}
		return mutationMsg{verb: "added", err: err}
	}
}

// tick schedules the next one-second header refresh.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// refreshHeader reads the running-timer and pending-outbox status in a command
// (off the render path) and delivers it as a headerMsg. Both reads are local
// files; any error degrades to "no timer"/"no pending ops" — the header is
// informational and must never fail the loop.
func refreshHeader() tea.Cmd {
	return func() tea.Msg {
		var info headerInfo
		if state, err := timer.Load(); err == nil {
			info.timerRunning = true
			info.ticket = state.TicketCode
			info.activity = state.ActivityName
			info.startedAt = state.StartedAt
		}
		if n, err := outbox.Count(); err == nil {
			info.outboxN = n
		}
		return headerMsg(info)
	}
}
