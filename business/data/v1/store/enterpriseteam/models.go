package enterpriseteam

import (
	"github.com/lib/pq"
	"time"
)

type Team struct {
	ID                string      `db:"id" json:"id"`
	Name              string      `db:"name" json:"name"`
	InvoicingEntityID string      `db:"invoicing_entity_id" json:"InvoicingEntityID"`
	EnterpriseID      string      `db:"enterprise_id" json:"EnterpriseID"`
	PolicyID          string      `db:"policy_id" json:"policyID"`
	Description       string      `db:"description" json:"Description"`
	PaymentMethod     string      `db:"payment_method" json:"PaymentMethod"`
	UpdatedAt         time.Time   `db:"updated_at" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt         pq.NullTime `db:"deleted_at" json:"-"`
}

//NewTeamDTO represents the needed data to create a team
type NewTeamDTO struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	PaymentMethod     string `json:"paymentMethod"`
	EnterpriseID      string `json:"enterpriseId"`
	PolicyID          string `json:"policyId"`
	InvoicingEntityID string `json:"invoicingEntityId"`
}

// UpdateTeamDTO defines what information may be provided to modify an existing
// team. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateTeamDTO struct {
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	PaymentMethod     *string `json:"paymentMethod"`
	InvoicingEntityID *string `json:"invoicingEntityId"`
}
