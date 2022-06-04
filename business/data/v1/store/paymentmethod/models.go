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
