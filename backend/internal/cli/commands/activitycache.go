package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lockw1n/time-logger/internal/cli/api"
	"github.com/lockw1n/time-logger/internal/cli/config"
)

// The activities cache exists solely to make `tl add` work while the backend is
// unreachable. Logging an entry needs the activity *id*, but the user supplies a
// *name* — normally resolved by fetching /api/activities, which is exactly the
// call that fails offline. So every successful fetch is cached per company, and
// an offline lookup falls back to that cache. Without it, offline add could never
// build a payload to queue; with it, the queue-on-unreachable path is reachable.

const cacheDirName = "cache"

// activitiesForLookup returns the company's activities, preferring a live fetch
// (which refreshes the cache) and falling back to the cache only when the backend
// is unreachable. Any other error is returned as-is — a real answer, not offline.
func activitiesForLookup(ctx context.Context, client api.Client, companyID uint64) ([]api.Activity, error) {
	activities, err := client.Activities(ctx, companyID)
	if err == nil {
		saveActivitiesCache(companyID, activities) // best effort; never blocks the command
		return activities, nil
	}
	if errors.Is(err, api.ErrUnreachable) {
		if cached, cErr := loadActivitiesCache(companyID); cErr == nil && len(cached) > 0 {
			return cached, nil
		}
	}
	return nil, err
}

func activitiesCachePath(companyID uint64) (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, cacheDirName, "activities-"+strconv.FormatUint(companyID, 10)+".json"), nil
}

// saveActivitiesCache persists the activity list for offline reuse. Failures are
// swallowed: a missing cache only costs offline convenience, never correctness.
func saveActivitiesCache(companyID uint64, activities []api.Activity) {
	path, err := activitiesCachePath(companyID)
	if err != nil {
		return
	}
	data, err := json.Marshal(activities)
	if err != nil {
		return
	}
	_ = config.WriteFileAtomic(path, data, 0o600)
}

func loadActivitiesCache(companyID uint64) ([]api.Activity, error) {
	path, err := activitiesCachePath(companyID)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var activities []api.Activity
	if err := json.Unmarshal(data, &activities); err != nil {
		return nil, fmt.Errorf("parsing activities cache: %w", err)
	}
	return activities, nil
}
