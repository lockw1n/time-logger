// Package outbox implements the offline write queue and its replay engine. When
// the backend is unreachable, `tl add`/`tl stop` persist the intended entry as a
// JSON op under ~/.config/tl/outbox/; the next command that reaches the API
// replays them oldest-first (see sync.go).
//
// Only entry *creates* are queueable. Edits, deletes and invoices fail fast when
// offline: queueing a mutation against server state you cannot see is a
// corruption source (you'd be editing an id that may not exist, or deleting the
// wrong row after a concurrent change), not a convenience. Creates are deferred
// because the server's unique index + the 409 merge reconciles a duplicate cell
// — but note replay is at-least-once and NOT idempotent: a create redelivered
// after its file survived an accepted commit double-counts (see sync.go).
package outbox

import (
	"time"

	"github.com/lockw1n/time-logger/internal/cli/api"
)

// KindCreateEntry is the only queueable operation kind. It exists as a named
// constant so the on-disk schema is explicit and future kinds (if ever added)
// can be discriminated rather than assumed.
const KindCreateEntry = "create_entry"

// Op is one queued operation, serialized as a single JSON file in the outbox.
// The payload is the exact CreateEntryRequest body so replay reconstructs the
// original intent verbatim — no re-resolution of activities or dates at sync time.
type Op struct {
	Kind      string                 `json:"kind"`
	CreatedAt time.Time              `json:"created_at"`
	Payload   api.CreateEntryRequest `json:"payload"`
	// MergeOnConflict lets replay resolve a 409 by summing minutes into the
	// existing entry without prompting (the interactive merge decision the user
	// would have made is deferred to, and pre-authorized for, replay time).
	MergeOnConflict bool `json:"merge_on_conflict"`
	// Attempts and LastError are bookkeeping filled in as an op is retried or
	// moved to failed/, so `tl sync --list` can explain why an op is stuck.
	Attempts  int    `json:"attempts"`
	LastError string `json:"last_error"`
}
