package outbox

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ErrLocked means another tl process is already replaying the outbox. It is not
// a failure — the caller (background sync) should simply skip and let the holder
// finish; the foreground `tl sync` reports it and exits without replaying.
var ErrLocked = errors.New("outbox is locked by another sync")

// processAlive reports whether a process with the given pid is running. It is a
// package var so tests can simulate a dead holder without needing a real
// dead-but-unreaped pid, which is impossible to arrange portably.
var processAlive = func(pid int) bool {
	if pid <= 0 {
		return false
	}
	// On Unix, signal 0 performs error checking without sending a signal:
	// ESRCH means no such process, EPERM means it exists but we may not signal it.
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// acquireLock claims outbox/.lock for this process. A pre-existing lock held by
// a live process yields ErrLocked. A lock left by a dead process (crashed
// mid-sync) is stale: it is broken with a notice and the acquisition retried
// once. The returned release removes the lock.
func acquireLock(dir string, notice io.Writer) (release func(), err error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("creating %s: %w", dir, err)
	}
	path := filepath.Join(dir, lockName)

	switch err := publishLock(path); {
	case err == nil:
		return func() { _ = os.Remove(path) }, nil
	case !errors.Is(err, fs.ErrExist):
		return nil, fmt.Errorf("acquiring lock: %w", err)
	}

	// The lock exists — decide whether its holder is still alive.
	pid, ok := readLockPID(path)
	if ok && processAlive(pid) {
		return nil, ErrLocked
	}

	// Stale (dead or unreadable pid): break it and try once more. A second
	// ErrExist now means a live process grabbed it in the race window — honor it.
	fmt.Fprintf(notice, "breaking stale outbox lock (pid %d no longer running)\n", pid)
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("breaking stale lock: %w", err)
	}
	switch err := publishLock(path); {
	case err == nil:
		return func() { _ = os.Remove(path) }, nil
	case errors.Is(err, fs.ErrExist):
		return nil, ErrLocked
	default:
		return nil, fmt.Errorf("acquiring lock: %w", err)
	}
}

// publishLock atomically creates the lock file already containing this process's
// pid: it writes a temp file, then hardlinks it into place. os.Link fails with
// fs.ErrExist if the target exists (the exclusive-create guarantee) and, unlike
// O_EXCL followed by a separate write, never exposes an empty lock file that a
// racing reader could misjudge as stale.
func publishLock(path string) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".lock-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	fmt.Fprintf(tmp, "%d\n", os.Getpid())
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Link(tmp.Name(), path)
}

// readLockPID reads the pid recorded in a lock file. ok is false when the file
// is missing or its contents aren't an integer (an unreadable lock is treated as
// stale so a garbage file can't wedge the outbox forever).
func readLockPID(path string) (pid int, ok bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	pid, err = strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, false
	}
	return pid, true
}
