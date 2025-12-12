package models

import "time"

type InvoiceLine struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	InvoiceID uint64 `gorm:"column:invoice_id;not null"`
	EntryID   uint64 `gorm:"column:entry_id;not null"`

	Hours  float64 `gorm:"column:hours;not null"`
	Rate   float64 `gorm:"column:rate;not null"`
	Amount float64 `gorm:"column:amount;not null"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	Invoice *Invoice `gorm:"foreignKey:InvoiceID"`
	Entry   *Entry   `gorm:"foreignKey:EntryID"`
}

func (InvoiceLine) TableName() string {
	return "invoice_lines"
}
