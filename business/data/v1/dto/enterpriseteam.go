package dto

//NewTeam represents the needed data to create a team
type NewTeam struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	PaymentMethod     string `json:"paymentMethod"`
	EnterpriseID      string `json:"enterpriseId"`
	PolicyID          string `json:"policyId"`
	InvoicingEntityID string `json:"invoicingEntityId"`
}

// UpdateTeam defines what information may be provided to modify an existing
// team. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateTeam struct {
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	PaymentMethod     *string `json:"paymentMethod"`
	InvoicingEntityID *string `json:"invoicingEntityId"`
}
