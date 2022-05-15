package aggregator

import (
	"context"
	"fmt"
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

func (s Store) Create(ctx context.Context, na NewAggregator, now time.Time) (Aggregator, error) {
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
	INSERT INTO "public".aggregator
	(id, name, code, api_key, provider_timeout, active, type, payment_by_tgs, logo_url, updated_at, created_at, deleted_at)
	VALUES
	(:id, :name, :code, :api_key, :provider_timeout, :active, :type, :payment_by_tgs, :logo_url, :updated_at, :created_at, null)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, agr); err != nil {
		return Aggregator{}, fmt.Errorf("failed to create aggregator: %v", err)
	}

	return agr, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (Aggregator, error) {
	var agg Aggregator

	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT * FROM aggregator AS a WHERE a.id = :id
`
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &agg); err != nil {
		return agg, fmt.Errorf("aggregator %s not found", id)
	}

	return agg, nil
}
