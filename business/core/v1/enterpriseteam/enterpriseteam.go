package enterpriseteam

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprise"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterpriseteam"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/invoicingentity"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/teampolicy"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log             *zap.SugaredLogger
	db              *sqlx.DB
	teamStore       enterpriseteam.Store
	invoicingStore  invoicingentity.Store
	enterpriseStore enterprise.Store
	policyStore     teampolicy.Store
}

func NewCore(db *sqlx.DB, log *zap.SugaredLogger) Core {
	return Core{
		log:             log,
		db:              db,
		teamStore:       enterpriseteam.NewStore(db, log),
		invoicingStore:  invoicingentity.NewStore(db, log),
		enterpriseStore: enterprise.NewStore(db, log),
		policyStore:     teampolicy.NewStore(db, log),
	}
}

func (c Core) Create(ctx context.Context, nt dto.NewTeam, now time.Time) (enterpriseteam.Team, error) {
	//check if the invoicing entity, enterprise and the policy exist
	if _, err := c.invoicingStore.QueryByID(ctx, nt.InvoicingEntityID); err != nil {
		return enterpriseteam.Team{}, err
	}
	if _, err := c.enterpriseStore.QueryByID(ctx, nt.EnterpriseID); err != nil {
		return enterpriseteam.Team{}, err
	}
	if _, err := c.policyStore.QueryByID(ctx, nt.PolicyID); err != nil {
		return enterpriseteam.Team{}, err
	}

	t, err := c.teamStore.Create(ctx, nt, now)
	if err != nil {
		return enterpriseteam.Team{}, err
	}
	return t, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]enterpriseteam.Team, error) {
	ts, err := c.teamStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []enterpriseteam.Team{}, err
	}
	return ts, nil
}

func (c Core) QueryById(ctx context.Context, id string) (enterpriseteam.Team, error) {
	t, err := c.teamStore.QueryByID(ctx, id)
	if err != nil {
		return enterpriseteam.Team{}, err
	}
	return t, nil
}

func (c Core) QueryByEnterpriseID(ctx context.Context, id string) (enterpriseteam.Team, error) {
	t, err := c.teamStore.QueryByEnterpriseID(ctx, id)
	if err != nil {
		return enterpriseteam.Team{}, err
	}
	return t, nil
}

func (c Core) Update(ctx context.Context, id string, ut dto.UpdateTeam, now time.Time) (enterpriseteam.Team, error) {
	t, err := c.teamStore.QueryByID(ctx, id)
	if err != nil {
		return enterpriseteam.Team{}, err
	}
	if ut.Name != nil {
		t.Name = *ut.Name
	}
	if ut.InvoicingEntityID != nil {
		t.InvoicingEntityID = *ut.InvoicingEntityID
	}
	if ut.Description != nil {
		t.Description = *ut.Description
	}
	if ut.PaymentMethod != nil {
		t.PaymentMethod = *ut.PaymentMethod
	}
	if err := c.teamStore.Update(ctx, t, now); err != nil {
		return enterpriseteam.Team{}, err
	}
	return t, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (enterpriseteam.Team, error) {
	t, err := c.teamStore.QueryByID(ctx, id)
	if err != nil {
		return enterpriseteam.Team{}, err
	}
	if err := c.teamStore.Delete(ctx, id, now); err != nil {
		return enterpriseteam.Team{}, err
	}
	return t, nil
}
