package enterprisepolicy

import (
	"github.com/lib/pq"
	"time"
)

//TeamPolicy represents a policy apply to a team.
//A policy apply a set of rule to restrict the access
//to some offer to a collaborator.
//Currently, we support limitation by: hours
//It's also determine the paymentmethod type used by the collaborator
type TeamPolicy struct {
	ID                 string      `json:"id" db:"id"`
	Name               string      `json:"name" db:"name"`
	Description        string      `json:"description" db:"description"`
	CollaboratorBudget int         `json:"collaboratorBudget" db:"collaborator_budget"`
	StartServiceTime   string      `json:"startServiceTime" db:"start_service_time"`
	EndServiceTime     string      `json:"endServiceTime" db:"end_service_time"`
	BudgetType         string      `json:"budgetType" db:"budget_type"`
	EnterpriseID       string      `json:"enterpriseId" db:"enterprise_id"`
	BlockedDays        string      `json:"blockedDays" db:"blocked_days"` //@todo make it a []string, and make it updatable
	UpdatedAt          time.Time   `db:"updated_at" json:"-"`
	CreatedAt          time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt          pq.NullTime `db:"deleted_at" json:"-"`
}

//NewEnterprisePolicyDTO represents the required data to create a team policy
type NewEnterprisePolicyDTO struct {
	Name               string `json:"name" validate:"required"`
	Description        string `json:"description" validate:"required"`
	CollaboratorBudget int    `json:"collaboratorBudget" validate:"required"`
	StartServiceTime   string `json:"startServiceTime" validate:"required"`
	EndServiceTime     string `json:"endServiceTime" validate:"required"`
	BudgetType         string `json:"budgetType" validate:"required"`
	EnterpriseID       string `json:"enterpriseId" validate:"required"`
}

//UpdateEnterprisePolicyDTO defines what information may be provided to modify an existing
//// team policy. All fields are optional so clients can send just the fields they want
//// changed. It uses pointer fields so we can differentiate between a field that
//// was not provided and a field that was provided as explicitly blank. Normally
//// we do not want to use pointers to basic types but we make exceptions around
//// marshalling/unmarshalling.
type UpdateEnterprisePolicyDTO struct {
	Name               *string `json:"name"`
	Description        *string `json:"description"`
	CollaboratorBudget *int    `json:"collaboratorBudget"`
	StartServiceTime   *string `json:"startServiceTime"`
	EndServiceTime     *string `json:"endServiceTime"`
	BudgetType         *string `json:"budgetType"`
}
