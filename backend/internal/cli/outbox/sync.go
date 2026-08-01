package outbox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/entrylog"
)

// Report summarizes a replay run: how many ops were accepted by the server, how
// many were moved to failed/ this run, and how many remain pending afterwards.
type Report struct {
	Synced    int
	Failed    int
	Remaining int
}

// Options configures a replay run.
type Options struct {
	// Notice receives the stale-lock message; nil discards it.
	Notice io.Writer
	// Submit is the create-or-merge function, defaulting to entrylog.Submit. It
	// is injectable so tests can drive replay logic without a real merge path.
	Submit SubmitFunc
}

// SubmitFunc creates an entry, merging into an existing one on a 409, exactly as
// the interactive add/stop path does — but non-interactively.
type SubmitFunc func(ctx context.Context, client api.Client, req entrylog.Request, opts entrylog.Options) error

// Sync replays queued ops oldest-first, sequentially. Order matters: two ops for
// the same ticket/activity/day must merge deterministically, so they are never
// replayed concurrently. It holds an exclusive lock for the whole run; a second
// concurrent Sync returns ErrLocked without touching the queue.
//
// Per op the disposition is one of:
//   - success            → the op file is deleted.
//   - permanent failure  → the op is moved to failed/ with LastError; run continues.
//   - backend unreachable → the run stops immediately, leaving every remaining
//     op (including this one) untouched, to be retried next time.
//
// Delivery is at-least-once: an op whose create the server accepted but whose
// file was not deleted (a crash, or a failed unlink, in the small window between
// commit and remove) is replayed again. Replay is NOT idempotent — the 409
// handler SUMS minutes, so a redelivered create double-counts. Without a
// server-side idempotency key the 409 path cannot tell a redelivery apart from a
// genuine second entry for the same cell (which is meant to sum), so this window
// is an accepted trade-off of the schema-free outbox rather than something the
// client can fully close. A redelivery caused by a failed unlink is warned about
// loudly (below) instead of silently reconciled.
func Sync(ctx context.Context, client api.Client, opts Options) (Report, error) {
	notice := opts.Notice
	if notice == nil {
		notice = io.Discard
	}
	submit := opts.Submit
	if submit == nil {
		submit = entrylog.Submit
	}

	dir, err := Dir()
	if err != nil {
		return Report{}, err
	}

	release, err := acquireLock(dir, notice)
	if err != nil {
		return Report{}, err
	}
	defer release()

	entries, err := Pending()
	if err != nil {
		return Report{}, err
	}

	var rep Report
	for _, e := range entries {
		// Honor cancellation between ops: a mid-flight op that was cancelled
		// stays queued (at-least-once), and remaining ops are left untouched.
		if ctx.Err() != nil {
			break
		}

		// An op we can't understand (corrupt JSON, or an unknown future kind)
		// must never be replayed as a zeroed payload; quarantine the raw file so
		// its bytes survive for inspection, and move on.
		if e.Op.Kind != KindCreateEntry {
			if qErr := quarantineRaw(e.Filename); qErr != nil {
				return finalize(rep), qErr
			}
			rep.Failed++
			continue
		}

		err := replayOne(ctx, client, submit, e.Op)
		switch classify(err) {
		case dispSuccess:
			rep.Synced++
			if rmErr := remove(e.Filename); rmErr != nil {
				// The server accepted this op but its queue file survived. We
				// cannot delete it, so the next run will replay it and the
				// 409-merge will double-count. There is no safe automatic
				// recovery (re-queuing to failed/ needs the same unlink that just
				// failed), so warn loudly and keep going with the rest of the run.
				fmt.Fprintf(notice,
					"warning: synced %s but could not remove its queue file %s (%v) — it may be re-sent and double-counted on the next sync; delete it manually from the outbox directory\n",
					e.Op.Payload.TicketCode, e.Filename, rmErr)
			}
		case dispPermanent:
			if mvErr := moveToFailed(e.Filename, e.Op, err.Error()); mvErr != nil {
				return finalize(rep), mvErr
			}
			rep.Failed++
		case dispStop:
			// Unreachable / auth / transient server error: stop the run and
			// leave this op and all after it for the next attempt.
			return finalize(rep), nil
		}
	}

	return finalize(rep), nil
}

// finalize fills Remaining from the current on-disk pending count so the number
// reflects reality even after early stops or partial runs.
func finalize(rep Report) Report {
	if n, err := Count(); err == nil {
		rep.Remaining = n
	}
	return rep
}

// replayOne reconstructs the original submit request from the stored payload and
// runs it non-interactively. The activity name is intentionally left blank: the
// merge path matches on activity id (which the payload carries) and only uses the
// name for human output, which is discarded here.
func replayOne(ctx context.Context, client api.Client, submit SubmitFunc, op Op) error {
	req := entrylog.Request{
		CompanyID:  op.Payload.CompanyID,
		TicketCode: op.Payload.TicketCode,
		Activity:   api.Activity{ID: op.Payload.ActivityID},
		Date:       op.Payload.Date,
		Minutes:    op.Payload.DurationMinutes,
		Comment:    op.Payload.Comment,
	}
	return submit(ctx, client, req, entrylog.Options{
		Out:         io.Discard,
		SkipConfirm: true,
	})
}

// disposition is how the replay loop should treat an op's result.
type disposition int

const (
	dispSuccess   disposition = iota // accepted (created or merged) → delete
	dispPermanent                    // can never succeed → move to failed/
	dispStop                         // transient/offline → stop, retry later
)

// classify maps a submit error to a disposition. The key distinction is offline
// (ErrUnreachable) and other transient conditions (auth expiry, 5xx) — which
// should be retried, so the run stops — versus permanent rejections (validation
// 4xx, or the merge path's "no active contract"/limit errors) which will never
// succeed and must be quarantined rather than retried forever or silently dropped.
func classify(err error) disposition {
	if err == nil {
		return dispSuccess
	}
	// Offline, interrupted and auth failures are retryable: stop and leave the
	// queue intact. ErrInterrupted (cancel/timeout) is deliberately not
	// quarantined — the op stays for the next run (at-least-once).
	if errors.Is(err, api.ErrUnreachable) || errors.Is(err, api.ErrInterrupted) || errors.Is(err, api.ErrUnauthorized) {
		return dispStop
	}
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		if apiErr.Status >= http.StatusInternalServerError {
			return dispStop // 5xx is a transient server problem, retry later.
		}
		return dispPermanent // other 4xx (e.g. 400 validation) can never succeed.
	}
	// A plain error here comes from the merge path (no active contract for the
	// date, or the 24h/entry limit) — permanent for this op as written.
	return dispPermanent
}
