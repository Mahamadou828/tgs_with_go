package invoicingentity

import (
	"github.com/lib/pq"
	"time"
)

type InvoicingEntity struct {
	ID           string      `db:"id" json:"id"`
	Denomination string      `db:"denomination" json:"denomination"`
	EnterpriseID string      `db:"enterprise_id" json:"enterpriseId"`
	Number       string      `db:"number" json:"number"`
	Vat          string      `db:"vat" json:"vat"`
	Street       string      `db:"street" json:"street"`
	PostalCode   string      `db:"postal_code" json:"postalCode"`
	Town         string      `db:"town" json:"town"`
	Country      string      `db:"country" json:"country"`
	UpdatedAt    time.Time   `db:"updated_at" json:"-"`
	CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt    pq.NullTime `db:"deleted_at" json:"-"`
}
