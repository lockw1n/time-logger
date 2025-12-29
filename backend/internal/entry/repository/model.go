package repository

import "time"

type entryModel struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	ConsultantID    uint64    `gorm:"column:consultant_id;not null;uniqueIndex:ux_entry;index:idx_entries_consultant_company_date"`
	CompanyID       uint64    `gorm:"column:company_id;not null;uniqueIndex:ux_entry;index:idx_entries_consultant_company_date"`
	ContractID      uint64    `gorm:"column:contract_id;not null;uniqueIndex:ux_entry"`
	TicketID        uint64    `gorm:"column:ticket_id;not null;uniqueIndex:ux_entry"`
	ActivityID      uint64    `gorm:"column:activity_id;not null;uniqueIndex:ux_entry"`
	Date            time.Time `gorm:"column:date;type:date;not null;uniqueIndex:ux_entry;index:idx_entries_consultant_company_date"`
	DurationMinutes int       `gorm:"column:duration_minutes;not null"`
	Comment         *string   `gorm:"column:comment"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (entryModel) TableName() string {
	return "entries"
}
