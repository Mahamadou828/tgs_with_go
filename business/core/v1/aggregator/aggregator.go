package aggregator

import (
	"context"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	log      *zap.SugaredLogger
	db       *sqlx.DB
	aws      *aws.Client
	aggStore aggregator.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.Client) Core {
	return Core{
		log:      log,
		db:       db,
		aws:      aws,
		aggStore: aggregator.NewStore(log, db, aws),
	}
}

func (c Core) Create(ctx context.Context, na dto.NewAggregator, now time.Time) (aggregator.Aggregator, error) {
	agg, err := c.aggStore.Create(ctx, na, now)
	if err != nil {
		return aggregator.Aggregator{}, err
	}
	return agg, nil
}

func (c Core) Update(ctx context.Context, id string, ua dto.UpdateAggregator, now time.Time) (aggregator.Aggregator, error) {
	dbAgg, err := c.QueryByID(ctx, id)
	if err != nil {
		return aggregator.Aggregator{}, err
	}

	if ua.Code != nil {
		dbAgg.Code = *ua.Code
	}
	if ua.ApiKey != nil {
		dbAgg.ApiKey = *ua.ApiKey
	}
	if ua.Name != nil {
		dbAgg.Name = *ua.Name
	}
	if ua.ProviderTimeout != nil {
		dbAgg.ProviderTimeout = *ua.ProviderTimeout
	}
	if ua.Active != nil {
		dbAgg.Active = *ua.Active
	}
	if ua.Type != nil {
		dbAgg.Type = *ua.Type
	}
	if ua.PaymentByTGS != nil {
		dbAgg.PaymentByTGS = *ua.PaymentByTGS
	}
	if ua.LogoURL != nil {
		dbAgg.LogoURL = *ua.LogoURL
	}

	if err := c.aggStore.Update(ctx, id, dbAgg, now); err != nil {
		return aggregator.Aggregator{}, err
	}

	return dbAgg, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (aggregator.Aggregator, error) {
	agg, err := c.aggStore.QueryByID(ctx, id)
	if err != nil {
		return aggregator.Aggregator{}, err
	}
	return agg, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]aggregator.Aggregator, error) {
	return c.aggStore.Query(ctx, pageNumber, rowsPerPage)
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (aggregator.Aggregator, error) {
	agg, err := c.QueryByID(ctx, id)

	if err != nil {
		return aggregator.Aggregator{}, err
	}

	if err := c.aggStore.Delete(ctx, id, now); err != nil {
		return aggregator.Aggregator{}, err
	}

	return agg, err
}
