package models

import "time"

type Entry struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ConsultantID           uint64 `gorm:"column:consultant_id;not null"`
	CompanyID              uint64 `gorm:"column:company_id;not null"`
	ConsultantAssignmentID uint64 `gorm:"column:consultant_assignment_id;not null"`

	TicketID *uint64 `gorm:"column:ticket_id"`
	LabelID  *uint64 `gorm:"column:label_id"`

	Date            time.Time `gorm:"column:date;type:date;not null"`
	DurationMinutes int       `gorm:"column:duration_minutes;not null"`
	Comment         *string   `gorm:"column:comment"`

	HourlyRateSnapshot *float64 `gorm:"column:hourly_rate_snapshot"`
	CurrencySnapshot   *string  `gorm:"column:currency_snapshot"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Consultant        *Consultant           `gorm:"foreignKey:ConsultantID"`
	Company           *Company              `gorm:"foreignKey:CompanyID"`
	ConsultantCompany *ConsultantAssignment `gorm:"foreignKey:ConsultantAssignmentID"`
	Ticket            *Ticket               `gorm:"foreignKey:TicketID"`
	Label             *Label                `gorm:"foreignKey:LabelID"`
}

func (Entry) TableName() string {
	return "entries"
}
