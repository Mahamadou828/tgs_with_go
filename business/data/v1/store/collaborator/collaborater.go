package collaborator

import (
	"context"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	aws *aws.AWS
	db  *sqlx.DB
}

type CreateCollaboratorParams struct {
	Params   dto.NewCollaborator
	ApiKey   string
	AggID    string
	Budget   int
	StripeID string
	//AwsID is the aws cognito sub
	AwsID string
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Store {
	return Store{
		log: log,
		aws: aws,
		db:  db,
	}
}

func (s Store) Query(ctx context.Context, pageNumber int, rows int) ([]Collaborator, error) {
	data := struct {
		Offset int `db:"offset"`
		Rows   int `db:"rows_per_page"`
	}{
		Offset: (pageNumber - 1) * rows,
		Rows:   rows,
	}
	const q = `
	SELECT 
		id, 
		email, 
		phone_number, 
		name, 
		stripe_id, 
		api_key, 
		aggregator_id, 
		active, 
		cognito_id, 
		is_monthly_active, 
		is_cgu_accepted, 
		role,
		enterprise_team_id,
		enterprise_id,
		budget,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."user" 
	WHERE deleted_at IS NULL
	ORDER BY 
		id 
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY
`
	var collabs []Collaborator

	if err := database.NamedQuerySlice[Collaborator](ctx, s.log, s.db, q, data, &collabs); err != nil {
		return []Collaborator{}, err
	}

	if collabs == nil {
		return []Collaborator{}, nil
	}

	return collabs, nil
}

func (s Store) QueryByEnterprise(ctx context.Context, entID string, pageNumber, rows int) ([]Collaborator, error) {
	data := struct {
		Offset       int    `db:"Offset"`
		Rows         int    `db:"Rows"`
		EnterpriseID string `db:"enterprise_id"`
	}{
		Offset:       pageNumber,
		Rows:         rows,
		EnterpriseID: entID,
	}
	const q = `
	SELECT 
		id, 
		email, 
		phone_number, 
		name, 
		stripe_id, 
		api_key, 
		aggregator_id, 
		active, 
		cognito_id, 
		is_monthly_active, 
		is_cgu_accepted, 
		role,
		enterprise_team_id,
		enterprise_id,
		budget,
		updated_at, 
		created_at, 
		deleted_at  
	FROM 
		"public"."user" 
	WHERE 
		deleted_at IS NULL 
	AND 
		enterprise_id = :enterprise_id
	ORDER BY 
		id 
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY
`
	var collabs []Collaborator

	if err := database.NamedQuerySlice[Collaborator](ctx, s.log, s.db, q, data, &collabs); err != nil {
		return []Collaborator{}, err
	}

	if collabs == nil {
		return []Collaborator{}, nil
	}

	return collabs, nil
}

func (s Store) QueryByID(ctx context.Context, collabID string) (Collaborator, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: collabID,
	}
	const q = `
	SELECT 
		id, 
		email, 
		phone_number, 
		name, 
		stripe_id, 
		api_key, 
		aggregator_id, 
		active, 
		cognito_id, 
		is_monthly_active, 
		is_cgu_accepted, 
		role,
		enterprise_team_id,
		enterprise_id,
		budget,
		updated_at, 
		created_at, 
		deleted_at  
	FROM 
		"public"."user" 
	WHERE 
		deleted_at IS NULL 
	AND 
		id = :id
`
	var collab Collaborator

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &collab); err != nil {
		return Collaborator{}, err
	}

	return collab, nil
}

func (s Store) QueryByEmail(ctx context.Context, agg, email string) (Collaborator, error) {
	data := struct {
		Email        string `db:"email"`
		AggregatorID string `db:"aggregator_id"`
	}{
		Email:        email,
		AggregatorID: agg,
	}
	const q = `
	SELECT 
		id, 
		email, 
		phone_number, 
		name, 
		stripe_id, 
		api_key, 
		aggregator_id, 
		active, 
		cognito_id, 
		is_monthly_active, 
		is_cgu_accepted, 
		role,
		enterprise_team_id,
		enterprise_id,
		budget,
		updated_at, 
		created_at, 
		deleted_at  
	FROM 
		"public"."user" 
	WHERE 
		deleted_at IS NULL 
	AND 
		email 		= :email
	AND
		aggregator_id 	= :aggregator_id
`
	var collab Collaborator

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &collab); err != nil {
		return Collaborator{}, err
	}

	return collab, nil
}

func (s Store) Create(ctx context.Context, now time.Time, p CreateCollaboratorParams) (Collaborator, error) {
	co := Collaborator{
		ID:              validate.GenerateID(),
		Email:           p.Params.Email,
		PhoneNumber:     p.Params.PhoneNumber,
		Name:            p.Params.Name,
		StripeID:        p.StripeID,
		ApiKey:          p.ApiKey,
		AggregatorID:    p.AggID,
		EnterpriseID:    p.Params.EnterpriseID,
		TeamID:          p.Params.TeamID,
		Active:          p.Params.IsPhoneNumberVerified,
		CognitoID:       p.AwsID,
		IsMonthlyActive: false,
		IsCGUAccepted:   false,
		Role:            p.Params.Role,
		UpdatedAt:       now,
		CreatedAt:       now,
		Budget:          p.Budget,
		DeletedAt: pq.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	const q = `
	INSERT INTO "public"."user" 
		(id, aggregator_id, email, phone_number, name, budget, stripe_id, enterprise_id, enterprise_team_id, api_key, active, cognito_id, is_monthly_active, is_cgu_accepted, role, created_at, updated_at, deleted_at)
	VALUES
		(:id, :aggregator_id, :email, :phone_number, :name, :budget, :stripe_id, :enterprise_id, :enterprise_team_id, :api_key, :active, :cognito_id, false, :is_cgu_accepted, :role, :created_at, :updated_at, null)
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, co); err != nil {
		return Collaborator{}, err
	}

	return co, nil
}

func (s Store) Delete(ctx context.Context, id string, now time.Time) error {
	data := struct {
		ID        string      `db:"id"`
		DeletedAt pq.NullTime `db:"deleted_at"`
	}{
		ID: id,
		DeletedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	const q = `
	UPDATE 
		"public"."user" 
	SET 
		deleted_at = :deleted_at 
	WHERE 
		id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}

func (s Store) Update(ctx context.Context, id string, co Collaborator, now time.Time) error {
	data := struct {
		UpdatedAt        pq.NullTime `db:"updated_at"`
		ID               string      `db:"id"`
		Email            string      `db:"email"`
		PhoneNumber      string      `db:"phone_number"`
		Name             string      `db:"name"`
		Active           bool        `db:"active"`
		IsMonthlyActive  bool        `db:"is_monthly_active" json:"isMonthlyActive"`
		IsCGUAccepted    bool        `db:"is_cgu_accepted" json:"isCGUAccepted"`
		Role             string      `db:"role" json:"role"`
		EnterpriseTeamID string      `db:"enterprise_team_id"`
		Budget           int         `db:"budget"`
	}{
		UpdatedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
		ID:               id,
		Name:             co.Name,
		Email:            co.Email,
		PhoneNumber:      co.ApiKey,
		Active:           co.Active,
		Role:             co.Role,
		IsMonthlyActive:  co.IsMonthlyActive,
		IsCGUAccepted:    co.IsCGUAccepted,
		EnterpriseTeamID: co.TeamID,
		Budget:           co.Budget,
	}

	const q = `
		UPDATE
			"public"."user" 
		SET 
			name              	= :name,
			email             	= :email,
			phone_number      	= :phone_number,
			active            	= :active,
			role              	= :role,
			is_monthly_active 	= :is_monthly_active,
			is_cgu_accepted   	= :is_cgu_accepted,
			enterprise_team_id 	= :enterprise_team_id,
			budget   			= :budget,
			updated_at        	= :updated_at
		WHERE 
			id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
