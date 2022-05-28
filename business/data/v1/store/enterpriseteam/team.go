package enterpriseteam

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

func (s Store) Create(ctx context.Context, nt dto.NewTeam, now time.Time) (Team, error) {
	team := Team{
		ID:                validate.GenerateID(),
		Name:              nt.Name,
		InvoicingEntityID: nt.InvoicingEntityID,
		EnterpriseID:      nt.EnterpriseID,
		Description:       nt.Description,
		PaymentMethod:     nt.PaymentMethod,
		UpdatedAt:         now,
		PolicyID:          nt.PolicyID,
		CreatedAt:         now,
		DeletedAt: pq.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	const q = `
	INSERT INTO "public"."enterprise_team"
		(id, name, invoicing_entity_id, policy_id, enterprise_id, description, payment_method, created_at, updated_at, deleted_at) 
	VALUES 
		(:id, :name, :invoicing_entity_id, :policy_id, :enterprise_id, :description, :payment_method, :created_at, :updated_at, :deleted_at)
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, team); err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (Team, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		*
	FROM 
		"public"."enterprise_team" 
	WHERE deleted_at IS NULL AND id = :id 
`

	var team Team

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Store) QueryByEnterpriseID(ctx context.Context, id string) (Team, error) {
	data := struct {
		EnterpriseID string `db:"enterprise_id"`
	}{
		EnterpriseID: id,
	}

	const q = `
	SELECT 
		*
	FROM 
		"public"."enterprise_team" 
	WHERE deleted_at IS NULL AND enterprise_id = :enterprise_id 
`

	var team Team

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return Team{}, err
	}

	return team, nil
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]Team, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT 
		* 
	FROM 
		"public"."enterprise_team" 
	WHERE deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var teams []Team

	if err := database.NamedQuerySlice[Team](ctx, s.log, s.db, q, data, &teams); err != nil {
		return []Team{}, err
	}

	if teams == nil {
		return []Team{}, nil
	}

	return teams, nil
}
func (s Store) Update(ctx context.Context, t Team, now time.Time) error {
	data := struct {
		ID                string    `db:"id"`
		Name              string    `db:"name"`
		InvoicingEntityID string    `db:"invoicing_entity_id"`
		Description       string    `db:"description"`
		PaymentMethod     string    `db:"payment_method"`
		UpdatedAt         time.Time `db:"updated_at"`
	}{
		ID:                t.ID,
		Name:              t.Name,
		InvoicingEntityID: t.InvoicingEntityID,
		Description:       t.Description,
		PaymentMethod:     t.PaymentMethod,
		UpdatedAt:         now,
	}

	const q = `
	UPDATE 
		"public"."enterprise_team"
	SET 	
		name 				= :name
		invoicing_entity_id = :invoicing_entity_id
		description 		= :description
		payment_method		= :payment_method
		updated_at 			= :updated_at
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
	UPDATE "public"."enterprise_team" SET deleted_at = :deleted_at WHERE id = :id
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
