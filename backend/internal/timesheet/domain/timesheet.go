package domain

import (
	"time"

	activitydomain "github.com/lockw1n/time-logger/internal/activity/domain"
	ticketdomain "github.com/lockw1n/time-logger/internal/ticket/domain"
)

type Timesheet struct {
	ConsultantID uint64
	CompanyID    uint64
	Start        time.Time
	End          time.Time
	Rows         []TimesheetRow
	Totals       TimesheetTotals
}

type TimesheetRow struct {
	Ticket        ticketdomain.Ticket
	Activity      activitydomain.Activity
	Entries       []TimesheetEntry
	PerDayMinutes map[string]int
	TotalMinutes  int
}

type TimesheetEntry struct {
	ID              uint64
	Date            time.Time
	DurationMinutes int
	Comment         *string
}

type TimesheetTotals struct {
	PerDayMinutes  map[string]int
	OverallMinutes int
}
