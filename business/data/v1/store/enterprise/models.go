package enterprise

import (
	"github.com/lib/pq"
	"time"
)

type Enterprise struct {
	ID                string      `json:"id"`
	Code              string      `db:"code" json:"code"`
	Name              string      `db:"name" json:"name"`
	ContactEmail      string      `db:"contact_email" json:"contactEmail"`
	Description       string      `db:"description" json:"description"`
	LogoURL           string      `db:"logo_url" json:"logoUrl"`
	MaxCarbonEmission int         `db:"max_carbon_emission" json:"maxCarbonEmission"`
	Active            bool        `db:"active" json:"active"`
	UpdatedAt         time.Time   `db:"updated_at" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt         pq.NullTime `db:"deleted_at" json:"-"`
	PackID            string      `db:"pack_id" json:"packId"`
}