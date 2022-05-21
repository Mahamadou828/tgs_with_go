//Package user provide an api to interact with the business
//logic related to user CRUD and management
package user

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/user"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log       *zap.SugaredLogger
	db        *sqlx.DB
	aws       *aws.AWS
	userStore user.Store
	aggStore  aggregator.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Core {
	return Core{
		log:       log,
		db:        db,
		aws:       aws,
		userStore: user.NewStore(log, db, aws),
		aggStore:  aggregator.NewStore(log, db, aws),
	}
}

func (c Core) Create(ctx context.Context, aggregatorCode string, nu user.NewUser, now time.Time) (user.User, error) {
	agg, err := c.aggStore.QueryByID(ctx, aggregatorCode)

	if err != nil {
		return user.User{}, err
	}

	usr, err := c.userStore.Create(ctx, agg.ID, agg.ApiKey, nu, now)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Query(ctx context.Context, pages, rows int) ([]user.User, error) {
	usr, err := c.userStore.Query(ctx, pages, rows)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (user.User, error) {
	usr, err := c.userStore.QueryById(ctx, id)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Update(ctx context.Context, id string) (user.User, error) {
	usr, err := c.userStore.QueryById(ctx, id)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Delete(ctx context.Context, userId string, now time.Time) (user.User, error) {
	usr, err := c.userStore.Delete(ctx, userId, now)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}
