package dto

//NewInvoicingEntity represents the needed information to create a new Invoicing entity
type NewInvoicingEntity struct {
	Denomination string `json:"denomination" validate:"required"`
	EnterpriseID string `json:"enterpriseId" validate:"required"`
	Number       string `json:"number" validate:"required"`
	Vat          string `json:"vat" validate:"required"`
	Street       string `json:"street" validate:"required"`
	PostalCode   string `json:"postalCode" validate:"required"`
	Town         string `json:"town" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

// UpdateInvoicingEntity defines what information may be provided to modify an existing
// invoicing entity. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateInvoicingEntity struct {
	Denomination *string `json:"denomination"`
	Number       *string `json:"number"`
	Vat          *string `json:"vat"`
	Street       *string `json:"street"`
	PostalCode   *string `json:"postalCode"`
	Town         *string `json:"town"`
	Country      *string `json:"country"`
}
