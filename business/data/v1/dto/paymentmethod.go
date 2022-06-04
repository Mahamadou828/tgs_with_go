package dto

type NewPaymentMethod struct {
	Number      string `json:"number" validate:"required"`
	CVC         string `json:"cvc" validate:"required"`
	ExpireMonth string `json:"expireMonth" validate:"required"`
	ExpireYear  string `json:"expireYear" validate:"required"`
	Name        string `json:"name" validate:"required"`
	IsFavorite  bool   `json:"isFavorite" validate:"required"`
	Type        string `json:"type" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
	ReturnURL   string `json:"returnUrl" validate:"required"`
}

// UpdatePaymentMethod defines what information may be provided to modify an existing
// update paymentmethod method. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdatePaymentMethod struct {
	Name       *string `json:"name"`
	IsFavorite *bool   `json:"isFavorite"`
}
