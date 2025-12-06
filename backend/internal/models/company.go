package models

import "time"

// Company holds billing details used in invoices.
type Company struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	UID          string    `json:"uid"`
	AddressLine1 string    `json:"address_line1"`
	Zip          string    `json:"zip"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	Payment      string    `json:"payment_condition"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
