package domain

import "time"

type Period struct {
	Month string // YYYY-MM
	Start time.Time
	End   time.Time
}
