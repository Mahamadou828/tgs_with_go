package aggregator

import (
	"context"
	"fmt"
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

func (s Store) Create(ctx context.Context, na dto.NewAggregator, now time.Time) (Aggregator, error) {
	agr := Aggregator{
		ID:              validate.GenerateID(),
		Name:            na.Name,
		Code:            na.Code,
		ApiKey:          na.ApiKey,
		ProviderTimeout: na.ProviderTimeout,
		Active:          na.Active,
		Type:            na.Type,
		PaymentByTGS:    na.PaymentByTGS,
		LogoURL:         na.LogoURL,
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	const q = `
	INSERT INTO "public"."aggregator"
	(id, name, code, api_key, provider_timeout, active, type, payment_by_tgs, logo_url, updated_at, created_at, deleted_at)
	VALUES
	(:id, :name, :code, :api_key, :provider_timeout, :active, :type, :payment_by_tgs, :logo_url, :updated_at, :created_at, null)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, agr); err != nil {
		return Aggregator{}, fmt.Errorf("failed to create aggregator: %v", err)
	}

	return agr, nil
}

func (s Store) Update(ctx context.Context, id string, ua Aggregator, now time.Time) error {
	data := struct {
		UpdatedAt       pq.NullTime `db:"updated_at"`
		ID              string      `db:"id"`
		Name            string      `db:"name"`
		Code            string      `db:"code"`
		ApiKey          string      `db:"api_key"`
		ProviderTimeout int         `db:"provider_timeout"`
		Active          bool        `db:"active"`
		Type            string      `db:"type"`
		PaymentByTGS    bool        `db:"payment_by_tgs"`
		LogoURL         string      `db:"logo_url"`
	}{
		UpdatedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
		ID:              id,
		Name:            ua.Name,
		Code:            ua.Code,
		ApiKey:          ua.ApiKey,
		ProviderTimeout: ua.ProviderTimeout,
		Active:          ua.Active,
		Type:            ua.Type,
		PaymentByTGS:    ua.PaymentByTGS,
		LogoURL:         ua.LogoURL,
	}

	const q = `
	UPDATE 
		"public"."aggregator" 
	SET 
		name            	= :name,
		code            	= :code,
		api_key          	= :api_key,
		provider_timeout 	= :provider_timeout,
		active          	= :active,
		type            	= :type,
		payment_by_tgs    	= :payment_by_tgs,
		logo_url         	= :logo_url,
		updated_at          = :updated_at
	WHERE 
		id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}

func (s Store) QueryByID(ctx context.Context, id string) (Aggregator, error) {
	var agg Aggregator

	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT * FROM "public"."aggregator" AS a WHERE a.id = :id AND deleted_at IS NULL
`
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &agg); err != nil {
		return agg, fmt.Errorf("aggregator %s not found", id)
	}

	return agg, nil
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]Aggregator, error) {
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
		"public"."aggregator" 
	WHERE deleted_at IS NOT NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var aggrs []Aggregator

	if err := database.NamedQuerySlice[Aggregator](ctx, s.log, s.db, q, data, &aggrs); err != nil {
		return []Aggregator{}, err
	}

	if aggrs == nil {
		return []Aggregator{}, nil
	}

	return aggrs, nil
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
	UPDATE "public"."aggregator" SET deleted_at = :deleted_at WHERE id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
