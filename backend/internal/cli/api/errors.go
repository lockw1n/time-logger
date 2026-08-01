package api

import "errors"

var (
	// ErrUnauthorized signals a 401 on an authenticated call; it drives exit code 2.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrUnreachable signals that the request never reached the backend
	// (connection refused, DNS failure, no route). The offline outbox relies on
	// this: only an unreachable write is safe to queue and replay, because it is
	// certain the server did not act on it.
	ErrUnreachable = errors.New("backend unreachable")

	// ErrInterrupted signals a request whose outcome is unknown — it was
	// cancelled (Ctrl-C) or timed out after being sent, so the server may or may
	// not have committed it. Such a write must NOT be auto-queued: replaying a
	// request the server already applied would double-count via the merge path.
	ErrInterrupted = errors.New("request interrupted")
)

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
