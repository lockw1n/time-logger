package models

import "time"

// Entry represents a time log record for a specific ticket and date.
type Entry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Ticket    string    `json:"ticket" gorm:"index:idx_ticket;index:idx_ticket_date,priority:1;uniqueIndex:uniq_ticket_date,priority:1"`
	Label     string    `json:"label"`
	Hours     float64   `json:"hours" gorm:"type:numeric(6,2)"`                                                                 // quarter-hour precision enforced at DB level
	Date      time.Time `json:"date" gorm:"type:date;index:idx_ticket_date,priority:2;uniqueIndex:uniq_ticket_date,priority:2"` // stored as date-only (UTC-normalized)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
