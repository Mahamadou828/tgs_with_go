package invoicingentity

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprise"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/invoicingentity"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log             *zap.SugaredLogger
	db              *sqlx.DB
	invoicingStore  invoicingentity.Store
	enterpriseStore enterprise.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:             log,
		db:              db,
		invoicingStore:  invoicingentity.NewStore(db, log),
		enterpriseStore: enterprise.NewStore(db, log),
	}
}

func (c Core) Create(ctx context.Context, ni dto.NewInvoicingEntity, now time.Time) (invoicingentity.InvoicingEntity, error) {
	if _, err := c.enterpriseStore.QueryByID(ctx, ni.EnterpriseID); err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	i, err := c.invoicingStore.Create(ctx, ni, now)
	if err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	return i, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]invoicingentity.InvoicingEntity, error) {
	is, err := c.invoicingStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []invoicingentity.InvoicingEntity{}, err
	}
	return is, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (invoicingentity.InvoicingEntity, error) {
	i, err := c.invoicingStore.QueryByID(ctx, id)
	if err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	return i, nil
}

func (c Core) QueryByEnterpriseID(ctx context.Context, id string) (invoicingentity.InvoicingEntity, error) {
	i, err := c.invoicingStore.QueryByEnterpriseID(ctx, id)
	if err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	return i, nil
}

func (c Core) Update(ctx context.Context, id string, ni dto.UpdateInvoicingEntity, now time.Time) (invoicingentity.InvoicingEntity, error) {
	i, err := c.invoicingStore.QueryByID(ctx, id)
	if err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}

	if ni.Street != nil {
		i.Street = *ni.Street
	}
	if ni.Town != nil {
		i.Town = *ni.Town
	}
	if ni.Denomination != nil {
		i.Denomination = *ni.Denomination
	}
	if ni.PostalCode != nil {
		i.PostalCode = *ni.PostalCode
	}
	if ni.Country != nil {
		i.Country = *ni.Country
	}
	if ni.Town != nil {
		i.Town = *ni.Town
	}
	if ni.Number != nil {
		i.Number = *ni.Number
	}
	if err := c.invoicingStore.Update(ctx, i, now); err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	return i, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (invoicingentity.InvoicingEntity, error) {
	i, err := c.invoicingStore.QueryByID(ctx, id)
	if err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	if err := c.invoicingStore.Delete(ctx, id, now); err != nil {
		return invoicingentity.InvoicingEntity{}, err
	}
	return i, nil
}
