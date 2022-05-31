package collaborator

import (
	"github.com/lib/pq"
	"time"
)

//Collaborator represents a new collaborator attach to a enterprise
type Collaborator struct {
	ID              string      `db:"id" json:"id"`
	Email           string      `db:"email" json:"email"`
	PhoneNumber     string      `db:"phone_number" json:"phoneNumber"`
	Name            string      `db:"name" json:"name"`
	StripeID        string      `db:"stripe_id" json:"stripeId"`
	ApiKey          string      `db:"api_key" json:"apiKey"`
	AggregatorID    string      `db:"aggregator_id" json:"aggregatorId"`
	EnterpriseID    string      `db:"enterprise_id" json:"enterpriseId"`
	TeamID          string      `db:"enterprise_team_id," json:"teamId"`
	Active          bool        `db:"active" json:"active"`
	CognitoID       string      `db:"cognito_id" json:"cognitoId"`
	IsMonthlyActive bool        `db:"is_monthly_active" json:"isMonthlyActive"`
	IsCGUAccepted   bool        `db:"is_cgu_accepted" json:"isCGUAccepted"`
	Budget          int         `db:"budget" json:"budget"`
	Role            string      `db:"role" json:"role"`
	UpdatedAt       time.Time   `db:"updated_at" json:"-"`
	CreatedAt       time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt       pq.NullTime `db:"deleted_at" json:"-"`
}
