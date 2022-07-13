package user

import (
	"github.com/lib/pq"
	"time"
)

//User represents a new user create under the aggregator logic
type User struct {
	ID              string         `db:"id" json:"id"`
	Email           string         `db:"email" json:"email"`
	PhoneNumber     string         `db:"phone_number" json:"phoneNumber"`
	Name            string         `db:"name" json:"name"`
	StripeID        string         `db:"stripe_id" json:"stripeId"`
	ApiKey          string         `db:"api_key" json:"apiKey"`
	AggregatorID    string         `db:"aggregator_id" json:"aggregatorId"`
	Active          bool           `db:"active" json:"active"`
	CognitoID       string         `db:"cognito_id" json:"cognitoId"`
	IsMonthlyActive bool           `db:"is_monthly_active" json:"isMonthlyActive"`
	IsCGUAccepted   bool           `db:"is_cgu_accepted" json:"isCGUAccepted"`
	Role            pq.StringArray `db:"role" json:"role"`
	UpdatedAt       time.Time      `db:"updated_at" json:"-"`
	CreatedAt       time.Time      `db:"created_at" json:"createdAt"`
	DeletedAt       pq.NullTime    `db:"deleted_at" json:"-"`
}

//NewUserDTO contains the minimal needed information to create a new user
type NewUserDTO struct {
	Email                 string   `json:"email" validate:"required,email"`
	PhoneNumber           string   `json:"phoneNumber" validate:"required"`
	Name                  string   `json:"name" validate:"required"`
	IsPhoneNumberVerified bool     `json:"isPhoneNumberVerified" validate:"required"`
	IsCGUAccepted         bool     `json:"isCGUAccepted" validate:"required"`
	Role                  []string `json:"role" validate:"required,oneof=USER"`
	Password              string   `json:"password" validate:"required"`
}

// UpdateUserDTO defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateUserDTO struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	PhoneNumber     *string  `json:"phoneNumber"`
	Active          *bool    `json:"active"`
	IsMonthlyActive *bool    `json:"isMonthlyActive"`
	Role            []string `json:"role"`
	IsCGUAccepted   *bool    `json:"isCGUAccepted"`
}

type LoginDTO struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type RefreshTokenDTO struct {
	ID           string `json:"id" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type VerifyConfirmationCodeDTO struct {
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type ConfirmNewPasswordDTO struct {
	NewPassword string `json:"new_password" validate:"required"`
	Code        string `json:"code" validate:"required"`
	ID          string `json:"id" validate:"required"`
}
