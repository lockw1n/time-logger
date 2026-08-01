package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
	"github.com/lockw1n/time-logger/internal/cli/outbox"
	"github.com/lockw1n/time-logger/internal/cli/parse"
)

// queuedError marks that a create was accepted into the offline outbox instead
// of being confirmed by the backend. It carries exit code 3 (per the overview:
// "operation queued offline") and has already printed its own message, so main
// prints nothing extra — same contract as exitError.
type queuedError struct{}

func (queuedError) Error() string { return "queued for offline sync" }
func (queuedError) ExitCode() int { return 3 }

// isQueued reports whether err is the queued-offline signal, so callers like
// `tl stop` can treat "queued" as a successful capture (clear the timer) rather
// than a failure to preserve state.
func isQueued(err error) bool {
	var q queuedError
	return errors.As(err, &q)
}

// submitEntry runs the shared create-or-merge flow and, when — and only when —
// the *initial create* never reached the backend, queues it for later sync
// instead of failing. A failure from the merge path (where an update may already
// have committed) is returned unchanged even if it is unreachable, because
// re-queuing it as a fresh create could double-count; entrylog signals the
// safe-to-queue case with ErrCreateUnreachable. Any other error (validation,
// auth, conflict, a declined merge) is likewise returned unchanged. On a
// successful queue it returns queuedError (exit 3).
func submitEntry(ctx context.Context, client api.Client, req entrylog.Request, opts entrylog.Options) error {
	err := entrylog.Submit(ctx, client, req, opts)
	if !errors.Is(err, entrylog.ErrCreateUnreachable) {
		return err
	}

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
		// If we cannot even persist the op, the work really would be lost —
		// surface the original unreachable error joined with the write failure.
		return fmt.Errorf("%w (also failed to queue for sync: %v)", err, qErr)
	}

	pending, _ := outbox.Count()
	fmt.Fprintf(queueOut(opts.Out), "api unreachable — queued %s on %s for sync (%d pending)\n",
		parse.FormatMinutes(req.Minutes), req.TicketCode, pending)
	return queuedError{}
}

// queueOut falls back to a discarding writer if the caller passed no output sink,
// so the queued message never panics on a nil writer.
func queueOut(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}
