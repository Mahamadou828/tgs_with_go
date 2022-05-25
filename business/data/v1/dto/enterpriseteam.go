package dto

//NewTeam represents the needed data to create a team
type NewTeam struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	PaymentMethod     string `json:"paymentMethod"`
	EnterpriseID      string `json:"enterpriseId"`
	InvoicingEntityID string `json:"invoicingEntityId"`
}
