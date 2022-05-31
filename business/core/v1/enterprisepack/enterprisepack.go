package enterprisepack

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprisepack"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log       *zap.SugaredLogger
	db        *sqlx.DB
	packStore enterprisepack.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:       log,
		db:        db,
		packStore: enterprisepack.NewStore(db, log),
	}
}

func (c Core) Create(ctx context.Context, np dto.NewPack, now time.Time) (enterprisepack.EnterprisePack, error) {
	p, err := c.packStore.Create(ctx, np, now)
	if err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	return p, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]enterprisepack.EnterprisePack, error) {
	ps, err := c.packStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []enterprisepack.EnterprisePack{}, err
	}
	return ps, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (enterprisepack.EnterprisePack, error) {
	p, err := c.packStore.QueryByID(ctx, id)
	if err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	return p, nil
}

func (c Core) Update(ctx context.Context, id string, up dto.UpdatePack, now time.Time) (enterprisepack.EnterprisePack, error) {
	p, err := c.packStore.QueryByID(ctx, id)
	if err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	if up.Name != nil {
		p.Name = *up.Name
	}
	if up.MaxActiveCollaboratorPerMonth != nil {
		p.MaxActiveCollaboratorPerMonth = *up.MaxActiveCollaboratorPerMonth
	}
	if up.SendExpenseReport != nil {
		p.SendExpenseReport = *up.SendExpenseReport
	}
	if up.CanCustomizeReport != nil {
		p.CanCustomizeReport = *up.CanCustomizeReport
	}
	if up.SendMonthlyReport != nil {
		p.SendMonthlyReport = *up.SendMonthlyReport
	}
	if up.CanCustomizeReport != nil {
		p.CanCustomizeReport = *up.CanCustomizeReport
	}

	if err := c.packStore.Update(ctx, p, now); err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	return p, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (enterprisepack.EnterprisePack, error) {
	p, err := c.packStore.QueryByID(ctx, id)
	if err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	if err := c.packStore.Delete(ctx, id, now); err != nil {
		return enterprisepack.EnterprisePack{}, err
	}
	return p, nil
}
