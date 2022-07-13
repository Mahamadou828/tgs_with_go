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

//NewInvoicingEntityDTO represents the needed information to create a new Invoicing entity
type NewInvoicingEntityDTO struct {
	Denomination string `json:"denomination" validate:"required"`
	EnterpriseID string `json:"enterpriseId" validate:"required"`
	Number       string `json:"number" validate:"required"`
	Vat          string `json:"vat" validate:"required"`
	Street       string `json:"street" validate:"required"`
	PostalCode   string `json:"postalCode" validate:"required"`
	Town         string `json:"town" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

// UpdateInvoicingEntityDTO defines what information may be provided to modify an existing
// invoicing entity. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateInvoicingEntityDTO struct {
	Denomination *string `json:"denomination"`
	Number       *string `json:"number"`
	Vat          *string `json:"vat"`
	Street       *string `json:"street"`
	PostalCode   *string `json:"postalCode"`
	Town         *string `json:"town"`
	Country      *string `json:"country"`
}
