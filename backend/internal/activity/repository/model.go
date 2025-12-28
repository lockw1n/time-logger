package repository

import "time"

type activityModel struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	CompanyID uint64  `gorm:"column:company_id;not null"`
	Name      string  `gorm:"column:name;not null"`
	Color     *string `gorm:"column:color"`
	Billable  bool    `gorm:"column:billable;default:true"`
	Priority  int     `gorm:"column:priority;default:1"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (activityModel) TableName() string {
	return "activities"
}
