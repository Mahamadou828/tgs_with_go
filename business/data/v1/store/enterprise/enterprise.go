package enterprise

import (
	"context"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

func (s Store) Create(ctx context.Context, nent dto.NewEnterprise, now time.Time) (Enterprise, error) {
	en := Enterprise{
		ID:                validate.GenerateID(),
		Code:              validate.GenerateEnterpriseCode(),
		Name:              nent.Name,
		ContactEmail:      nent.ContactEmail,
		Description:       nent.Description,
		LogoURL:           nent.LogoURL,
		MaxCarbonEmission: nent.MaxCarbonEmission,
		Active:            true,
		UpdatedAt:         now,
		CreatedAt:         now,
		PackID:            nent.PackID,
		DeletedAt: pq.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	const q = `
	INSERT INTO "public"."enterprise" 
		(id, pack_id, code, name, contact_email, description, logo_url, max_carbon_emission, active, created_at, updated_at, deleted_at)	
	VALUES 
		(:id, :pack_id, :code, :name, :contact_email, :description, :logo_url, :max_carbon_emission, :active, :created_at, :updated_at, :deleted_at)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, en); err != nil {
		return Enterprise{}, err
	}

	return en, nil
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]Enterprise, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      pageNumber,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT 
		id,
		code, 
		name, 
		contact_email, 
		description, 
		logo_url, 
		max_carbon_emission,
		active,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise" 
	WHERE deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var ents []Enterprise

	if err := database.NamedQuerySlice[Enterprise](ctx, s.log, s.db, q, data, &ents); err != nil {
		return []Enterprise{}, err
	}

	if ents == nil {
		return []Enterprise{}, nil
	}

	return ents, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (Enterprise, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	var ent Enterprise

	const q = `
	SELECT 
		id,
		code, 
		name, 
		contact_email, 
		description, 
		logo_url, 
		max_carbon_emission,
		active,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."enterprise"
	WHERE 
		id = :id 
	AND 
		deleted_at IS NULL
`

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &ent); err != nil {
		return Enterprise{}, err
	}

	return ent, nil
}

func (s Store) QueryByCode(ctx context.Context, code string) (Enterprise, error) {
	data := struct {
		Code string `db:"code"`
	}{
		Code: code,
	}

	var ent Enterprise

	const q = `
	SELECT 
		id,
		code, 
		name, 
		contact_email, 
		description, 
		logo_url, 
		max_carbon_emission,
		active,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."enterprise"
	WHERE 
		code = :code 
	AND 
		deleted_at IS NULL
`

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &ent); err != nil {
		return Enterprise{}, err
	}

	return ent, nil
}

func (s Store) Update(ctx context.Context, en Enterprise, now time.Time) error {
	data := struct {
		ID                string    `db:"id"`
		Code              string    `db:"code"`
		PackID            string    `db:"pack_id"`
		Name              string    `db:"name"`
		ContactEmail      string    `db:"contact_email"`
		Description       string    `db:"description"`
		LogoURL           string    `db:"logo_url"`
		MaxCarbonEmission int       `db:"max_carbon_emission"`
		UpdatedAt         time.Time `db:"updated_at"`
	}{
		ID:                en.ID,
		Code:              en.Code,
		PackID:            en.PackID,
		Name:              en.Name,
		ContactEmail:      en.ContactEmail,
		Description:       en.Description,
		LogoURL:           en.LogoURL,
		MaxCarbonEmission: en.MaxCarbonEmission,
		UpdatedAt:         now,
	}

	const q = `
		UPDATE
			"public"."enterprise" 
		SET 
			name              	= :name,
			contact_email     	= :contact_email,
			pack_id      	  	= :pack_id,
			code      	  		= :code,
			description         = :description,
			logo_url            = :logo_url,
			max_carbon_emission = :max_carbon_emission,
			updated_at       	= :updated_at
		WHERE 
			id = :id
`

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
		"public"."enterprise"
	SET 
		deleted_at = :deleted_at
	WHERE 	
		id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}
	return nil
}
