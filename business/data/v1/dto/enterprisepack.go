package dto

//NewPack represents minimal needed info to create a pack
type NewPack struct {
	Name                          string `json:"name"`
	SendMonthlyReport             bool   `json:"sendMonthlyReport"`
	CanCustomizeReport            bool   `json:"canCustomizeReport"`
	SendExpenseReport             bool   `json:"sendExpenseReport"`
	MaxActiveCollaboratorPerMonth int    `json:"maxActiveCollaboratorPerMonth"`
}
