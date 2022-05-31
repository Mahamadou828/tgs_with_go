package teampolicy

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprise"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/teampolicy"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log             *zap.SugaredLogger
	db              *sqlx.DB
	enterpriseStore enterprise.Store
	policyStore     teampolicy.Store
}

func NewCore(db *sqlx.DB, log *zap.SugaredLogger) Core {
	return Core{
		log:             log,
		db:              db,
		enterpriseStore: enterprise.NewStore(db, log),
		policyStore:     teampolicy.NewStore(db, log),
	}
}

func (c Core) Create(ctx context.Context, nt dto.NewTeamPolicy, now time.Time) (teampolicy.TeamPolicy, error) {
	//check if the invoicing entity, enterprise and the policy exist
	if _, err := c.enterpriseStore.QueryByID(ctx, nt.EnterpriseID); err != nil {
		return teampolicy.TeamPolicy{}, err
	}

	t, err := c.policyStore.Create(ctx, nt, now)
	if err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	return t, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]teampolicy.TeamPolicy, error) {
	ts, err := c.policyStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []teampolicy.TeamPolicy{}, err
	}
	return ts, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (teampolicy.TeamPolicy, error) {
	t, err := c.policyStore.QueryByID(ctx, id)
	if err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	return t, nil
}

func (c Core) QueryByEnterpriseID(ctx context.Context, id string) (teampolicy.TeamPolicy, error) {
	t, err := c.policyStore.QueryByEnterpriseID(ctx, id)
	if err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	return t, nil
}

func (c Core) Update(ctx context.Context, id string, ut dto.UpdateTeamPolicy, now time.Time) (teampolicy.TeamPolicy, error) {
	t, err := c.policyStore.QueryByID(ctx, id)
	if err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	if ut.Name != nil {
		t.Name = *ut.Name
	}
	if ut.CollaboratorBudget != nil {
		t.CollaboratorBudget = *ut.CollaboratorBudget
	}
	if ut.Description != nil {
		t.Description = *ut.Description
	}
	if ut.BudgetType != nil {
		t.BudgetType = *ut.BudgetType
	}
	if ut.EndServiceTime != nil {
		t.EndServiceTime = *ut.EndServiceTime
	}
	if ut.StartServiceTime != nil {
		t.StartServiceTime = *ut.StartServiceTime
	}

	if err := c.policyStore.Update(ctx, id, t, now); err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	return t, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (teampolicy.TeamPolicy, error) {
	t, err := c.policyStore.QueryByID(ctx, id)
	if err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	if err := c.policyStore.Delete(ctx, id, now); err != nil {
		return teampolicy.TeamPolicy{}, err
	}
	return t, nil
}
