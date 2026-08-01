package outbox

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/lockw1n/time-logger/internal/cli/config"
)

const (
	dirName    = "outbox"
	failedName = "failed"
	lockName   = ".lock"
)

// Dir returns the outbox directory within the tl config dir (TL_CONFIG_DIR aware).
func Dir() (string, error) {
	base, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, dirName), nil
}

// FailedDir returns the subdirectory holding ops that can never succeed.
func FailedDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, failedName), nil
}

// Entry pairs a queued op with the file it lives in, so callers can list ops and
// still address the underlying file (to delete, move, or discard it).
type Entry struct {
	Filename string
	Op       Op
}

// Enqueue writes op to a new outbox file and returns its filename. The name is
// <unix-nanos>-<random4>.json: the nanosecond prefix makes a lexical sort equal
// a chronological one (fixed-width until year 2262), which is the replay order;
// the random suffix avoids collisions when two ops are queued in the same tick.
func Enqueue(op Op) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(op, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding op: %w", err)
	}

	name := fmt.Sprintf("%d-%s.json", time.Now().UnixNano(), rand4())
	if err := config.WriteFileAtomic(filepath.Join(dir, name), data, 0o600); err != nil {
		return "", err
	}
	return name, nil
}

func rand4() string {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		// A weak suffix is acceptable — the nanosecond prefix already provides
		// near-certain uniqueness; the random part only breaks same-tick ties.
		return "0000"
	}
	return hex.EncodeToString(b[:])
}

// Pending lists queued ops (not yet synced, not failed) in replay order.
func Pending() ([]Entry, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	return readEntries(dir)
}

// Failed lists ops that were moved aside because they can never succeed.
func Failed() ([]Entry, error) {
	dir, err := FailedDir()
	if err != nil {
		return nil, err
	}
	return readEntries(dir)
}

// Count reports how many ops are pending replay — the cheap check every command
// uses to decide whether to bother launching a background sync.
func Count() (int, error) {
	names, err := pendingFiles()
	if err != nil {
		return 0, err
	}
	return len(names), nil
}

// readEntries reads and parses every *.json op file directly in dir (never
// recursing into subdirectories), sorted by filename = chronological order. A
// file that fails to parse is surfaced as an Entry with a zero Op and a
// LastError note rather than aborting the whole listing.
func readEntries(dir string) ([]Entry, error) {
	names, err := jsonFiles(dir)
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(names))
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if errors.Is(err, fs.ErrNotExist) {
			// The file vanished between listing and reading — e.g. a concurrent
			// sync in another process just replayed and deleted it. Skip it rather
			// than failing the whole listing.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", name, err)
		}
		var op Op
		if err := json.Unmarshal(data, &op); err != nil {
			// Leave Kind empty so the replay engine can tell this op apart from a
			// valid one and quarantine the original bytes instead of replaying zeros.
			op = Op{LastError: fmt.Sprintf("unparseable op file: %v", err)}
		}
		entries = append(entries, Entry{Filename: name, Op: op})
	}
	return entries, nil
}

// pendingFiles returns the sorted names of the pending op files.
func pendingFiles() ([]string, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	return jsonFiles(dir)
}

// jsonFiles returns the sorted *.json filenames directly inside dir. Directory
// entries (e.g. failed/) and the .lock file are ignored. A missing dir is not an
// error — an empty outbox is the normal case.
func jsonFiles(dir string) ([]string, error) {
	des, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	var names []string
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		if filepath.Ext(de.Name()) != ".json" {
			continue
		}
		names = append(names, de.Name())
	}
	sort.Strings(names)
	return names, nil
}

// remove deletes a pending op file after it has been synced.
func remove(filename string) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, filename)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", filename, err)
	}
	return nil
}

// moveToFailed records why op failed and moves its file into failed/. It writes
// the annotated op to failed/ first, then removes the pending file, so a crash
// mid-move never loses the op (at worst it exists in both places briefly, and the
// pending copy would simply be retried).
func moveToFailed(filename string, op Op, cause string) error {
	failedDir, err := FailedDir()
	if err != nil {
		return err
	}
	op.LastError = cause
	op.Attempts++

	data, err := json.MarshalIndent(op, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding failed op: %w", err)
	}
	if err := config.WriteFileAtomic(filepath.Join(failedDir, filename), data, 0o600); err != nil {
		return err
	}
	return remove(filename)
}

// quarantineRaw moves an op file into failed/ byte-for-byte, without parsing or
// re-encoding it. It exists for files the replay engine cannot understand
// (corrupt JSON, an unknown kind): re-marshaling them would destroy the original
// content that a human needs to diagnose the problem.
func quarantineRaw(filename string) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	failedDir, err := FailedDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(failedDir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", failedDir, err)
	}
	if err := os.Rename(filepath.Join(dir, filename), filepath.Join(failedDir, filename)); err != nil {
		return fmt.Errorf("quarantining %s: %w", filename, err)
	}
	return nil
}

// Discard permanently deletes a failed op by filename. It only ever touches the
// failed/ directory (never a live pending op) and rejects any path components in
// the name so a caller-supplied string can't escape the directory.
func Discard(filename string) error {
	if filename != filepath.Base(filename) || filename == "." || filename == ".." {
		return fmt.Errorf("invalid op filename %q", filename)
	}
	failedDir, err := FailedDir()
	if err != nil {
		return err
	}
	path := filepath.Join(failedDir, filename)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("no failed op named %q", filename)
		}
		return fmt.Errorf("removing %s: %w", filename, err)
	}
	return nil
}
