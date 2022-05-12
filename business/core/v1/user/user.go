//Package user provide an api to interact with the business
//logic related to user CRUD and management
package user

import (
	"context"
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
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Core {
	return Core{
		log:       log,
		db:        db,
		aws:       aws,
		userStore: user.NewStore(log, db, aws),
	}
}

func (c Core) Create(ctx context.Context, aggregatorID, apiKey string, nu user.NewUser, now time.Time) (user.User, error) {
	usr, err := c.userStore.Create(ctx, aggregatorID, apiKey, nu, now)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}