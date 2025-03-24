package models

import "time"

type UserPayment struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	PaymentMethod    string    `json:"payment_method"`
	TransactionID    *string   `json:"transaction_id"` // Nullable
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentStatus    string    `json:"payment_status"`
	PaymentDate      time.Time `json:"payment_date"`
	BillingAddressID *int      `json:"billing_address_id"` // Nullable
	Description      *string   `json:"description"`        // Nullable
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedAt        time.Time `json:"created_at"`
}
