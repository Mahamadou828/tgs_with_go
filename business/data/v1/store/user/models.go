package user

import "time"

//User represents a new user create under the aggregator logic
type User struct {
	ID              string    `db:"id" json:"id"`
	Email           string    `db:"email" json:"email"`
	PhoneNumber     string    `db:"phoneNumber" json:"phoneNumber"`
	Name            string    `db:"name" json:"name"`
	StripeID        string    `db:"stripeId" json:"-"`
	ApiKey          string    `db:"apiKey" json:"-"`
	AggregatorID    string    `db:"aggregatorId" json:"-"`
	Active          bool      `db:"active" json:"active"`
	CognitoID       string    `db:"cognitoId" json:"-"`
	IsMonthlyActive bool      `db:"isMonthlyActive" json:"-"`
	IsCGUAccepted   bool      `db:"isCGUAccepted" json:"isCGUAccepted"`
	Role            []string  `db:"role" json:"-"`
	UpdatedAt       time.Time `db:"updatedAt" json:"-"`
	CreatedAt       time.Time `db:"createdAt" json:"createdAt"`
	DeletedAt       time.Time `db:"deletedAt" json:"-"`
}

//NewUser contains the minimal needed information to create a new user
type NewUser struct {
	Email                 string   `json:"email" validate:"required,email"`
	PhoneNumber           string   `json:"phoneNumber" validate:"required"`
	Name                  string   `json:"name" validate:"required"`
	IsPhoneNumberVerified bool     `json:"isPhoneNumberVerified" validate:"required"`
	IsCGUAccepted         bool     `json:"isCGUAccepted" validate:"required"`
	Role                  []string `json:"role" validate:"required"`
	Password              string   `json:"password" validate:"required"`
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	Email         *string `json:"email" validate:"omitempty,email"`
	PhoneNumber   *string `json:"phoneNumber"`
	Name          *string `json:"name"`
	IsCGUAccepted *bool   `json:"isCGUAccepted"`
}
