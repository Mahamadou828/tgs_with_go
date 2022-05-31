package dto

//NewUser contains the minimal needed information to create a new user
type NewUser struct {
	Email                 string `json:"email" validate:"required,email"`
	PhoneNumber           string `json:"phoneNumber" validate:"required"`
	Name                  string `json:"name" validate:"required"`
	IsPhoneNumberVerified bool   `json:"isPhoneNumberVerified" validate:"required"`
	IsCGUAccepted         bool   `json:"isCGUAccepted" validate:"required"`
	Role                  string `json:"role" validate:"required,oneof=USER"`
	Password              string `json:"password" validate:"required"`
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	Name            *string `json:"name"`
	Email           *string `json:"email" validate:"omitempty,email"`
	PhoneNumber     *string `json:"phoneNumber"`
	Active          *bool   `json:"active"`
	IsMonthlyActive *bool   `json:"isMonthlyActive"`
	Role            *string `json:"role"`
	IsCGUAccepted   *bool   `json:"isCGUAccepted"`
}

type Login struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type RefreshToken struct {
	ID           string `json:"id" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type VerifyConfirmationCode struct {
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type ConfirmNewPassword struct {
	NewPassword string `json:"new_password" validate:"required"`
	Code        string `json:"code" validate:"required"`
	ID          string `json:"id" validate:"required"`
}
