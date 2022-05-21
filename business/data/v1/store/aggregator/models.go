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
