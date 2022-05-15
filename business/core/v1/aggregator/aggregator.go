package aggregator

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	log      *zap.SugaredLogger
	db       *sqlx.DB
	aws      *aws.AWS
	aggStore aggregator.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Core {
	return Core{
		log:      log,
		db:       db,
		aws:      aws,
		aggStore: aggregator.NewStore(log, db, aws),
	}
}

func (c Core) Create(ctx context.Context, na aggregator.NewAggregator, now time.Time) (aggregator.Aggregator, error) {
	if err := validate.Check(na); err != nil {
		return aggregator.Aggregator{}, err
	}

	agg, err := c.aggStore.Create(ctx, na, now)
	if err != nil {
		return aggregator.Aggregator{}, err
	}
	return agg, nil
}

func (c Core) QueryByID(ctx context.Context, code string) (aggregator.Aggregator, error) {
	agg, err := c.aggStore.QueryByID(ctx, code)
	if err != nil {
		return aggregator.Aggregator{}, err
	}
	return agg, nil
}
