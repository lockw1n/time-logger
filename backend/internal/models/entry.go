package models

import "time"

// Entry represents a time log record for a specific ticket and date.
type Entry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Ticket    string    `json:"ticket"`
	Label     string    `json:"label"`
	Hours     float64   `json:"hours"`
	Date      time.Time `json:"date" gorm:"type:timestamp"` // stores UTC
	CreatedAt time.Time `json:"created_at"`
}
