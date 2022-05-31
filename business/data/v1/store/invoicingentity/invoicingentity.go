package invoicingentity

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

func NewStore(db *sqlx.DB, log *zap.SugaredLogger) Store {
	return Store{
		log: log,
		db:  db,
	}
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]InvoicingEntity, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT 
		id, 
		denomination, 
		enterprise_id, 
		number, 
		vat, 
		street,
		postal_code,
		town, 
		country,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."invoicing_entity"
	WHERE 
		deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var is []InvoicingEntity

	if err := database.NamedQuerySlice[InvoicingEntity](ctx, s.log, s.db, q, data, &is); err != nil {
		return []InvoicingEntity{}, err
	}

	if is == nil {
		return []InvoicingEntity{}, nil
	}

	return is, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (InvoicingEntity, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		id, 
		denomination, 
		enterprise_id, 
		number, 
		vat, 
		street,
		postal_code,
		town, 
		country,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."invoicing_entity"
	WHERE
		deleted_at IS NULL 
	AND 
		id = :id 
`
	var i InvoicingEntity
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &i); err != nil {
		return InvoicingEntity{}, err
	}
	return i, nil
}

func (s Store) QueryByEnterpriseID(ctx context.Context, id string) (InvoicingEntity, error) {
	data := struct {
		EnterpriseID string `db:"enterprise_id"`
	}{
		EnterpriseID: id,
	}

	const q = `
	SELECT 
		id, 
		denomination, 
		enterprise_id, 
		number, 
		vat, 
		street,
		postal_code,
		town, 
		country,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."invoicing_entity"
	WHERE
		deleted_at IS NULL 
	AND 
		enterprise_id = :enterprise_id 
`
	var i InvoicingEntity
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &i); err != nil {
		return InvoicingEntity{}, err
	}
	return i, nil
}

func (s Store) Create(ctx context.Context, ni dto.NewInvoicingEntity, now time.Time) (InvoicingEntity, error) {
	i := InvoicingEntity{
		ID:           validate.GenerateID(),
		Denomination: ni.Denomination,
		EnterpriseID: ni.EnterpriseID,
		Number:       ni.Number,
		Vat:          ni.Vat,
		Street:       ni.Street,
		PostalCode:   ni.PostalCode,
		Town:         ni.Town,
		Country:      ni.Country,
		UpdatedAt:    now,
		CreatedAt:    now,
		DeletedAt: pq.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	const q = `
	INSERT INTO "public"."invoicing_entity"
		(id, enterprise_id, denomination, number, vat, street, postal_code, town, country, created_at, updated_at, deleted_at)
	VALUES
		(:id, :enterprise_id, :denomination, :number, :vat, :street, :postal_code, :town, :country, :created_at, :updated_at, :deleted_at)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, i); err != nil {
		return InvoicingEntity{}, err
	}

	return i, nil
}

func (s Store) Update(ctx context.Context, i InvoicingEntity, now time.Time) error {
	data := struct {
		ID           string    `db:"id"`
		Denomination string    `db:"denomination"`
		Number       string    `db:"number"`
		Vat          string    `db:"vat"`
		Street       string    `db:"street"`
		PostalCode   string    `db:"postal_code"`
		Town         string    `db:"town"`
		Country      string    `db:"country"`
		UpdatedAt    time.Time `db:"updated_at"`
	}{
		ID:           i.ID,
		Denomination: i.Denomination,
		Number:       i.Number,
		Vat:          i.Vat,
		Street:       i.Street,
		PostalCode:   i.PostalCode,
		Town:         i.Town,
		Country:      i.Country,
		UpdatedAt:    now,
	}

	const q = `
	UPDATE 
		"public"."invoicing_entity"
	SET 
		denomination 	= :denomination,
		number 			= :number,
		vat         	= :vat,
		street 			= :street,
		postal_code     = :postal_code,
		town 			= :town,
		country 		= :country,
		updated_at 		= :updated_at
	WHERE 	
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}

func (s Store) Delete(ctx context.Context, id string, now time.Time) error {
	data := struct {
		ID        string      `db:"id"`
		DeletedAt pq.NullTime `db:"deleted_at"`
	}{
		ID: id,
		DeletedAt: pq.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	const q = `
	UPDATE 
		"public"."invoicing_entity"
	SET
		deleted_at 		= :deleted_at
	WHERE 	
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
