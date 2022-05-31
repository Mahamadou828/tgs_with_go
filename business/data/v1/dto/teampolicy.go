package dto

//NewTeamPolicy represents the required data to create a team policy
type NewTeamPolicy struct {
	Name               string `json:"name" validate:"required"`
	Description        string `json:"description" validate:"required"`
	CollaboratorBudget int    `json:"collaboratorBudget" validate:"required"`
	StartServiceTime   string `json:"startServiceTime" validate:"required"`
	EndServiceTime     string `json:"endServiceTime" validate:"required"`
	BudgetType         string `json:"budgetType" validate:"required"`
	EnterpriseID       string `json:"enterpriseId" validate:"required"`
}

//UpdateTeamPolicy defines what information may be provided to modify an existing
//// team policy. All fields are optional so clients can send just the fields they want
//// changed. It uses pointer fields so we can differentiate between a field that
//// was not provided and a field that was provided as explicitly blank. Normally
//// we do not want to use pointers to basic types but we make exceptions around
//// marshalling/unmarshalling.
type UpdateTeamPolicy struct {
	Name               *string `json:"name" validate:"required"`
	Description        *string `json:"description" validate:"required"`
	CollaboratorBudget *int    `json:"collaboratorBudget" validate:"required"`
	StartServiceTime   *string `json:"startServiceTime" validate:"required"`
	EndServiceTime     *string `json:"endServiceTime" validate:"required"`
	BudgetType         *string `json:"budgetType" validate:"required"`
}
