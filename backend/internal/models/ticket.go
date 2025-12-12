package models

import "time"

type Ticket struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	CompanyID uint64 `gorm:"column:company_id;not null"`

	Code        string  `gorm:"column:code;not null"`
	Label       string  `gorm:"column:label;not null"`
	Description *string `gorm:"column:description"`
	IsActive    bool    `gorm:"column:is_active;default:true"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Company *Company `gorm:"foreignKey:CompanyID"`
}

func (Ticket) TableName() string {
	return "tickets"
}
