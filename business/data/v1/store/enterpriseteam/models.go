package enterpriseteam

import (
	"github.com/lib/pq"
	"time"
)

type Team struct {
	ID                string `db:"id" json:"id"`
	Name              string
	InvoicingEntityID string      `db:"facturation_entity_id" json:"FacturationEntityID"`
	EnterpriseID      string      `db:"enterprise_id" json:"EnterpriseID"`
	Description       string      `db:"description" json:"Description"`
	PaymentMethod     string      `db:"payment_method" json:"PaymentMethod"`
	UpdatedAt         time.Time   `db:"updated_at" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt         pq.NullTime `db:"deleted_at" json:"-"`
}
