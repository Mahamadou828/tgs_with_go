package enterprisepack

import (
	"github.com/lib/pq"
	"time"
)

//EnterprisePack represents an available pack for enterprise
type EnterprisePack struct {
	ID                            string      `db:"id" json:"id"`
	Name                          string      `db:"name" json:"name"`
	SendMonthlyReport             bool        `db:"send_monthly_report" json:"sendMonthlyReport"`
	CanCustomizeReport            bool        `db:"can_customize_report" json:"canCustomizeReport"`
	SendExpenseReport             bool        `db:"send_expense_report" json:"sendExpenseReport"`
	MaxActiveCollaboratorPerMonth int         `db:"max_active_collaborator_per_month" json:"maxActiveCollaboratorPerMonth"`
	UpdatedAt                     time.Time   `db:"updated_at" json:"-"`
	CreatedAt                     time.Time   `db:"created_at" json:"createdAt"`
	DeletedAt                     pq.NullTime `db:"deleted_at" json:"-"`
}
