package user

import (
	"context"
	"fmt"
	userdto "github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/userroutes/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/service/v1/stripe"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (s Store) Create(ctx context.Context, agg aggregator.Aggregator, nu userdto.NewUser, now time.Time) (User, error) {
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
		AggID:       agg.ID,
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
		ApiKey:          agg.ApiKey,
		AggregatorID:    agg.ID,
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
	(id, aggregator_id, email, phone_number, name, stripe_id, api_key, active, cognito_id, is_monthly_active, is_cgu_accepted, role, created_at, updated_at, deleted_at)
	VALUES
	(:id, :aggregator_id, :email, :phone_number, :name, :stripe_id, :api_key, :active, :cognito_id, false, :is_cgu_accepted, :role, :created_at, :updated_at, null)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		//@todo we should rollback the user creating in cognito and stripe
		return User{}, fmt.Errorf("failed to create user %v: %v", usr.Email, err)
	}

	return usr, nil
}

func (s Store) Update(ctx context.Context, id string, u User, now time.Time) error {
	data := struct {
		UpdatedAt       pq.NullTime `db:"updated_at"`
		ID              string      `db:"id"`
		Email           string      `db:"email"`
		PhoneNumber     string      `db:"phone_number"`
		Name            string      `db:"name"`
		Active          bool        `db:"active"`
		IsMonthlyActive bool        `db:"is_monthly_active" json:"isMonthlyActive"`
		IsCGUAccepted   bool        `db:"is_cgu_accepted" json:"isCGUAccepted"`
		Role            string      `db:"role" json:"role"`
	}{
		UpdatedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
		ID:              id,
		Name:            u.Name,
		Email:           u.Email,
		PhoneNumber:     u.ApiKey,
		Active:          u.Active,
		Role:            u.Role,
		IsMonthlyActive: u.IsMonthlyActive,
		IsCGUAccepted:   u.IsCGUAccepted,
	}

	const q = `
		UPDATE
			"public"."user" 
		SET 
			name              = :name,
			email             = :email,
			phone_number      = :phone_number,
			active            = :active,
			role              = :role,
			is_monthly_active = :is_monthly_active,
			is_cgu_accepted   = :is_cgu_accepted,
			updated_at        = :updated_at
		WHERE 
			id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}

func (s Store) Delete(ctx context.Context, id string, now time.Time) (User, error) {
	u, err := s.QueryById(ctx, id)

	if err != nil {
		return User{}, err
	}
	data := struct {
		ID        string      `db:"id"`
		DeletedAt pq.NullTime `db:"deleted_at"`
	}{
		DeletedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
		ID: id,
	}
	const q = `
	UPDATE "public"."user" SET deleted_at = :deleted_at WHERE id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return User{}, err
	}

	return u, nil
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]User, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT 
		* 
	FROM 
		"public"."user" 
	WHERE deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var users []User

	if err := database.NamedQuerySlice[User](ctx, s.log, s.db, q, data, &users); err != nil {
		return []User{}, err
	}

	if users == nil {
		return []User{}, nil
	}

	return users, nil
}

func (s Store) QueryById(ctx context.Context, id string) (User, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	var u User

	const q = `
	SELECT * FROM "public"."user"  WHERE id = :id AND deleted_at IS NULL
`

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &u); err != nil {
		return u, fmt.Errorf("user #{%s} not found", id)
	}

	return u, nil
}

func (s Store) QueryByEmailAndAggregator(ctx context.Context, email, aggr string) (User, error) {
	data := struct {
		Email      string `db:"email"`
		Aggregator string `db:"aggregator_id"`
	}{
		Email:      email,
		Aggregator: aggr,
	}

	var u User

	const q = `
	SELECT * FROM "public"."user"  WHERE email = :email AND aggregator_id = :aggregator_id AND deleted_at IS NULL 
`

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &u); err != nil {
		return u, fmt.Errorf("user with email %s not found", email)
	}

	return u, nil
}

func (s Store) QueryByCognitoID(ctx context.Context, email, phoneNumber, aggregator string) (User, error) {
	cognitoID, err := s.aws.Cognito.GenerateSub(email, phoneNumber, aggregator)
	if err != nil {
		return User{}, err
	}
	data := struct {
		CognitoID string `db:"cognito_id"`
	}{
		CognitoID: cognitoID,
	}

	var u User

	const q = `
	SELECT * FROM "public"."user"  WHERE cognito_id = :cognito_id AND deleted_at IS NULL
`

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &u); err != nil {
		return u, fmt.Errorf("user with sub %s not found", cognitoID)
	}

	return u, nil
}
