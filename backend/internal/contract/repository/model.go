package repository

import "time"

type contractModel struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement"`
	ConsultantID uint64     `gorm:"column:consultant_id;not null"`
	CompanyID    uint64     `gorm:"column:company_id;not null"`
	HourlyRate   float64    `gorm:"column:hourly_rate;not null"`
	Currency     string     `gorm:"column:currency;not null"`
	OrderNumber  string     `gorm:"column:order_number;not null"`
	PaymentTerms *string    `gorm:"column:payment_terms"`
	StartDate    time.Time  `gorm:"column:start_date"`
	EndDate      *time.Time `gorm:"column:end_date"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (contractModel) TableName() string {
	return "contracts"
}
