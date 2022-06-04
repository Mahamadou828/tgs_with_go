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
