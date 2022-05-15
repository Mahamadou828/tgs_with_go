package user

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/service/v1/stripe"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
	aws *aws.AWS
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Store {
	return Store{
		log: log,
		db:  db,
		aws: aws,
	}
}

func (s Store) Create(ctx context.Context, aggID, apiKey string, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, err
	}

	stripeID, err := stripe.CreateUser(nu.Email, nu.PhoneNumber, nu.Name)

	if err != nil {
		return User{}, err
	}

	cognitoID, err := s.aws.Cognito.CreateUser(aws.CognitoUser{
		Email:       nu.Email,
		PhoneNumber: nu.PhoneNumber,
		Name:        nu.Name,
		AggID:       aggID,
		IsActive:    nu.IsPhoneNumberVerified,
		Password:    nu.Password,
	})

	if err != nil {
		return User{}, err
	}

	usr := User{
		ID:              validate.GenerateID(),
		Email:           nu.Email,
		PhoneNumber:     nu.PhoneNumber,
		Name:            nu.Name,
		StripeID:        stripeID,
		ApiKey:          apiKey,
		AggregatorID:    aggID,
		Active:          nu.IsPhoneNumberVerified,
		CognitoID:       cognitoID,
		IsMonthlyActive: false,
		IsCGUAccepted:   nu.IsCGUAccepted,
		Role:            nu.Role,
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	const q = `
	INSERT INTO "public"."user" 
	(id, aggregatorId, email, phoneNumber, name, stripeId, apiKey, active, cognitoId, isMonthlyActive, isCGUAccepted, role, createdAt, updatedAt, deletedAt)
	VALUES
	(:id, :aggregatorId, :email, :phoneNumber, :name, :stripeId, :apiKey, :active, :cognitoId, false, :isCGUAccepted, :role, :createdAt, :updatedAt, null)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		//@todo we should rollback the user creating in cognito and stripe
		return User{}, fmt.Errorf("failed to create user %v: %v", usr.Email, err)
	}

	return usr, nil
}

//func (s Store) Update(ctx context.Context, id string, uu UpdateUser, now time.Time) (User, error) {
//
//}
//
//func (s Store) Delete(ctx context.Context, id string, now time.Time) (User, error) {
//
//}
//
//func (s Store) Query(ctx context.Context, now time.Time) (User, error) {
//
//}
//
//func (s Store) QueryById(ctx context.Context, id string, now time.Time) (User, error) {
//
//}
//
//func (s Store) QueryByEmail(ctx context.Context, email string, now time.Time) (User, error) {
//
//}
