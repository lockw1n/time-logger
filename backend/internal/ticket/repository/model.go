package repository

import "time"

type ticketModel struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement"`
	CompanyID   uint64  `gorm:"column:company_id;not null;uniqueIndex:ux_company_code"`
	Code        string  `gorm:"column:code;not null;uniqueIndex:ux_company_code"`
	Title       *string `gorm:"column:title"`
	Label       *string `gorm:"column:label"`
	Description *string `gorm:"column:description"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ticketModel) TableName() string {
	return "tickets"
}
