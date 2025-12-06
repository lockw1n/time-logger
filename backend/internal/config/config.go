package config

import (
	"os"
	"strings"
)

var defaultLabels = []string{"feature", "bug", "meeting", "research"}

// AllowedLabels returns the configured labels (env: ALLOWED_LABELS as comma-separated list),
// falling back to defaults if unset or empty.
func AllowedLabels() []string {
	raw := os.Getenv("ALLOWED_LABELS")
	if strings.TrimSpace(raw) == "" {
		return defaultLabels
	}

	parts := strings.Split(raw, ",")
	labels := make([]string, 0, len(parts))
	for _, p := range parts {
		val := strings.ToLower(strings.TrimSpace(p))
		if val == "" {
			continue
		}
		labels = append(labels, val)
	}
	if len(labels) == 0 {
		return defaultLabels
	}
	return labels
}
