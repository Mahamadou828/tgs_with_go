package dto

//NewTeamPolicy represents the required data to create a team
type NewTeamPolicy struct {
	Name               string `json:"name" validate:"required"`
	Description        string `json:"description" validate:"required"`
	CollaboratorBudget int    `json:"collaboratorBudget" validate:"required"`
	StartServiceTime   string `json:"startServiceTime" validate:"required"`
	EndServiceTime     string `json:"endServiceTime" validate:"required"`
	BudgetType         string `json:"budgetType" validate:"required"`
	EnterpriseID       string `json:"enterpriseId" validate:"required"`
}
