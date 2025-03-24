package models

import "time"

type ShippingAddress struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 *string   `json:"address_line_2"` // Use pointer for nullable field
	City         string    `json:"city"`
	State        *string   `json:"state"` // Use pointer for nullable field
	PostalCode   string    `json:"postal_code"`
	Country      string    `json:"country"`
	IsDefault    bool      `json:"is_default"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}
