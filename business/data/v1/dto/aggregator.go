package dto

type NewAggregator struct {
	Name            string `validate:"required" json:"name"`
	Code            string `validate:"required" json:"code"`
	ApiKey          string `validate:"required" json:"apiKey"`
	ProviderTimeout int    `validate:"required" json:"providerTimeout"`
	Active          bool   `validate:"required" json:"active"`
	Type            string `validate:"required" json:"type"`
	PaymentByTGS    bool   `validate:"required" json:"paymentByTgs"`
	LogoURL         string `validate:"required" json:"logoUrl"`
}

type UpdateAggregator struct {
	Name            *string `json:"name"`
	Code            *string `json:"code"`
	ApiKey          *string `json:"apiKey"`
	ProviderTimeout *int    `json:"providerTimeout"`
	Active          *bool   `json:"active"`
	Type            *string `json:"type"`
	PaymentByTGS    *bool   `json:"paymentByTgs"`
	LogoURL         *string `json:"logoUrl"`
}
