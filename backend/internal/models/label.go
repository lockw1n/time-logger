package models

import "time"

type Label struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	CompanyID uint64  `gorm:"column:company_id;not null"`
	Name      string  `gorm:"column:name;not null"`
	Color     *string `gorm:"column:color"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Company *Company `gorm:"foreignKey:CompanyID"`
}

func (Label) TableName() string {
	return "labels"
}
