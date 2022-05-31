package dto

//NewEnterprise represents the minimal information needed to
//create a new enterprise. The code of the enterprise will be
//generated withing the creation workflow and by default the enterprise
//will be active
type NewEnterprise struct {
	Name               string `json:"name" validate:"required"`
	ContactEmail       string `json:"contactEmail" validate:"required,email"`
	Description        string `json:"description" validate:"required"`
	LogoURL            string `json:"logoUrl" validate:"required"`
	MaxCarbonEmission  int    `json:"maxCarbonEmission" validate:"required"`
	BlockedProvider    string `json:"blockedProvider" validate:"required"`
	BlockedProductType string `json:"blockedProductType" validate:"required"`
	PackID             string `json:"packId" validate:"required"`
}

type UpdateEnterprise struct {
	Name               *string `json:"name"`
	ContactEmail       *string `json:"contactEmail" validate:"omitempty,email"`
	Description        *string `json:"description"`
	LogoURL            *string `json:"logoUrl"`
	MaxCarbonEmission  *int    `json:"maxCarbonEmission"`
	BlockedProvider    *string `json:"blockedProvider"`
	BlockedProductType *string `json:"blockedProductType"`
	PackID             *string `json:"packId"`
}
