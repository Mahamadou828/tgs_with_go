package aggregator

import (
	"github.com/lib/pq"
	"time"
)

type Aggregator struct {
	ID              string      `db:"id" json:"id"`
	Name            string      `db:"name" json:"name"`
	Code            string      `db:"code" json:"code"`
	ApiKey          string      `db:"api_key" json:"apiKey"`
	ProviderTimeout int         `db:"provider_timeout" json:"providerTimeout"`
	Active          bool        `db:"active" json:"active"`
	Type            string      `db:"type" json:"type"`
	PaymentByTGS    bool        `db:"payment_by_tgs" json:"paymentByTgs"`
	LogoURL         string      `db:"logo_url" json:"logoUrl"`
	UpdatedAt       time.Time   `db:"updated_at" json:"-"`
	CreatedAt       time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt       pq.NullTime `db:"deleted_at" json:"-"`
}

type NewAggregator struct {
	Name            string `validate:"required" json:"name"`
	Code            string `validate:"required" json:"code"`
	ApiKey          string `validate:"required" json:"apiKey"`
	ProviderTimeout int    `validate:"required" json:"providerTimeout"`
	Active          bool   `validate:"required" json:"active"`
	Type            string `validate:"required" json:"type"`
	PaymentByTGS    bool   `validate:"required" json:"paymentByTgs"`
	LogoURL         string `validate:"required" json:"logoUrl"`
}

type UpdateAggregator struct {
	Name            *string `json:"name"`
	Code            *string `json:"code"`
	ApiKey          *string `json:"apiKey"`
	ProviderTimeout *int    `json:"providerTimeout"`
	Active          *bool   `json:"active"`
	Type            *string `json:"type"`
	PaymentByTGS    *bool   `json:"paymentByTgs"`
	LogoURL         *string `json:"logoUrl"`
}
