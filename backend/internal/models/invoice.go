package models

import "time"
import "gorm.io/datatypes"

type Invoice struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ConsultantID uint64 `gorm:"column:consultant_id;not null"`
	CompanyID    uint64 `gorm:"column:company_id;not null"`

	PeriodStart time.Time `gorm:"column:period_start;not null"`
	PeriodEnd   time.Time `gorm:"column:period_end;not null"`

	TotalAmount   *float64          `gorm:"column:total_amount"`
	Currency      string            `gorm:"column:currency;not null"`
	Status        string            `gorm:"column:status;not null;default:'draft'"`
	InvoiceNumber *string           `gorm:"column:invoice_number"`
	OrderNumber   *string           `gorm:"column:order_number"`
	Metadata      datatypes.JSONMap `gorm:"column:metadata;type:jsonb"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Consultant *Consultant   `gorm:"foreignKey:ConsultantID"`
	Company    *Company      `gorm:"foreignKey:CompanyID"`
	Entries    []InvoiceLine `gorm:"foreignKey:InvoiceID"`
}

func (Invoice) TableName() string {
	return "invoices"
}
