package dto

//NewPack represents minimal needed info to create a pack
type NewPack struct {
	Name                          string `json:"name" validate:"required"`
	SendMonthlyReport             bool   `json:"sendMonthlyReport" validate:"required"`
	CanCustomizeReport            bool   `json:"canCustomizeReport" validate:"required"`
	SendExpenseReport             bool   `json:"sendExpenseReport" validate:"required"`
	MaxActiveCollaboratorPerMonth int    `json:"maxActiveCollaboratorPerMonth" validate:"required"`
}

type UpdatePack struct {
	Name                          *string `json:"name"`
	SendMonthlyReport             *bool   `json:"sendMonthlyReport"`
	CanCustomizeReport            *bool   `json:"canCustomizeReport"`
	SendExpenseReport             *bool   `json:"sendExpenseReport"`
	MaxActiveCollaboratorPerMonth *int    `json:"maxActiveCollaboratorPerMonth"`
}
