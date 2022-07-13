package paymentmethod

import (
	"time"

	"github.com/lib/pq"
)

type PaymentMethod struct {
	ID                string      `db:"id" json:"id"`
	Name              string      `db:"name" json:"name"`
	UserID            string      `db:"user_id" json:"userId"`
	StripeID          string      `db:"stripe_id" json:"stripeId"`
	DisplayCreditCard string      `db:"display_credit_card" json:"displayCreditCard"`
	Type              string      `db:"type" json:"type"`
	IsFavorite        bool        `db:"is_favorite" json:"isFavorite"`
	UpdatedAt         time.Time   `db:"updated_at" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt         pq.NullTime `db:"deleted_at" json:"-"`
}

type NewPaymentMethodDTO struct {
	Number      string `json:"number" validate:"required"`
	CVC         string `json:"cvc" validate:"required"`
	ExpireMonth string `json:"expireMonth" validate:"required"`
	ExpireYear  string `json:"expireYear" validate:"required"`
	Name        string `json:"name" validate:"required"`
	IsFavorite  bool   `json:"isFavorite" validate:"required"`
	Type        string `json:"type" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
	ReturnURL   string `json:"returnUrl" validate:"required"`
}

// UpdatePaymentMethodDTO defines what information may be provided to modify an existing
// update paymentmethod method. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdatePaymentMethodDTO struct {
	Name       *string `json:"name"`
	IsFavorite *bool   `json:"isFavorite"`
}
