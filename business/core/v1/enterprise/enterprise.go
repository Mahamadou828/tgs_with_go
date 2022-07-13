package enterprise

import (
	"context"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprise"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprisepack"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	enterpriseStore enterprise.Store
	packStore       enterprisepack.Store
	log             *zap.SugaredLogger
	db              *sqlx.DB
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		enterpriseStore: enterprise.NewStore(db, log),
		packStore:       enterprisepack.NewStore(db, log),
		log:             log,
		db:              db,
	}
}

func (c Core) Create(ctx context.Context, ne enterprise.NewEnterpriseDTO, now time.Time) (enterprise.Enterprise, error) {
	if _, err := c.packStore.QueryByID(ctx, ne.PackID); err != nil {
		return enterprise.Enterprise{}, err
	}
	e, err := c.enterpriseStore.Create(ctx, ne, now)

	if err != nil {
		return enterprise.Enterprise{}, err
	}

	return e, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]enterprise.Enterprise, error) {
	ents, err := c.enterpriseStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []enterprise.Enterprise{}, err
	}
	return ents, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (enterprise.Enterprise, error) {
	ent, err := c.enterpriseStore.QueryByID(ctx, id)
	if err != nil {
		return enterprise.Enterprise{}, err
	}
	return ent, nil
}

func (c Core) QueryByCode(ctx context.Context, code string) (enterprise.Enterprise, error) {
	ent, err := c.enterpriseStore.QueryByCode(ctx, code)
	if err != nil {
		return enterprise.Enterprise{}, err
	}
	return ent, nil
}

func (c Core) Update(ctx context.Context, id string, ue enterprise.UpdateEnterpriseDTO, now time.Time) (enterprise.Enterprise, error) {
	ent, err := c.enterpriseStore.QueryByID(ctx, id)
	if err != nil {
		return enterprise.Enterprise{}, err
	}

	if ue.Name != nil {
		ent.Name = *ue.Name
	}

	if ue.PackID != nil {
		ent.PackID = *ue.PackID
	}

	if ue.LogoURL != nil {
		ent.LogoURL = *ue.LogoURL
	}

	if ue.ContactEmail != nil {
		ent.ContactEmail = *ue.ContactEmail
	}

	if ue.Description != nil {
		ent.Description = *ue.Description
	}

	if ue.MaxCarbonEmission != nil {
		ent.MaxCarbonEmission = *ue.MaxCarbonEmission
	}

	if err := c.enterpriseStore.Update(ctx, ent, now); err != nil {
		return enterprise.Enterprise{}, err
	}

	return ent, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (enterprise.Enterprise, error) {
	ent, err := c.enterpriseStore.QueryByID(ctx, id)
	if err != nil {
		return enterprise.Enterprise{}, err
	}
	if err := c.enterpriseStore.Delete(ctx, id, now); err != nil {
		return enterprise.Enterprise{}, err
	}
	return ent, nil
}
