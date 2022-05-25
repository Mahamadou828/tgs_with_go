package collaborator

import (
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/collaborator"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	collaboraterStore collaborator.Store
	aws               *aws.AWS
	db                *sqlx.DB
	log               *zap.SugaredLogger
}

func NewCore(aws *aws.AWS, db *sqlx.DB, log *zap.SugaredLogger) Core {
	return Core{
		collaboraterStore: collaborator.NewStore(),
		aws:               aws,
		db:                db,
		log:               log,
	}
}

func (c Core) Login() {

}

func (c Core) RefreshToken() {

}

func (c Core) Create() {

}

func (c Core) QueryByID() {

}

func (c Core) Query() {

}

func (c Core) QueryByEnterprise() {

}

func (c Core) Update() {

}

func (c Core) Delete() {

}
