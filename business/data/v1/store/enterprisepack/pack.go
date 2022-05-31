package enterprisepack

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

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]EnterprisePack, error) {
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
		name, 
		send_monthly_report, 
		can_customize_report, 
		send_expense_report,
		max_active_collaborator_per_month,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise_pack"
	WHERE deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY
`
	var packs []EnterprisePack

	if err := database.NamedQuerySlice[EnterprisePack](ctx, s.log, s.db, q, data, &packs); err != nil {
		return []EnterprisePack{}, err
	}

	if packs == nil {
		return []EnterprisePack{}, nil
	}

	return packs, nil

}

func (s Store) Create(ctx context.Context, np dto.NewPack, now time.Time) (EnterprisePack, error) {
	pack := EnterprisePack{
		ID:                            validate.GenerateID(),
		Name:                          np.Name,
		SendMonthlyReport:             np.SendMonthlyReport,
		CanCustomizeReport:            np.CanCustomizeReport,
		SendExpenseReport:             np.SendExpenseReport,
		MaxActiveCollaboratorPerMonth: np.MaxActiveCollaboratorPerMonth,
		UpdatedAt:                     now,
		CreatedAt:                     now,
		DeletedAt: pq.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	const q = `
	INSERT INTO "public"."enterprise_pack"
		(id, name, send_monthly_report, can_customize_report, send_expense_report, max_active_collaborator_per_month)
	VALUES
		(:id, :name, :send_monthly_report, :can_customize_report, :send_expense_report, :max_active_collaborator_per_month)
`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, pack); err != nil {
		return EnterprisePack{}, err
	}

	return pack, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (EnterprisePack, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		id,
		name, 
		send_monthly_report, 
		can_customize_report, 
		send_expense_report,
		max_active_collaborator_per_month,
		updated_at, 
		created_at, 
		deleted_at 
	FROM 
		"public"."enterprise_pack"
	WHERE 
		deleted_at IS NOT NULL
	AND 
		id = :id
`

	var pack EnterprisePack

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &pack); err != nil {
		return EnterprisePack{}, err
	}

	return pack, nil
}

func (s Store) Update(ctx context.Context, p EnterprisePack, now time.Time) error {
	data := struct {
		ID                            string    `db:"id"`
		Name                          string    `db:"name"`
		SendMonthlyReport             bool      `db:"send_monthly_report"`
		CanCustomizeReport            bool      `db:"can_customize_report"`
		SendExpenseReport             bool      `db:"send_expense_report"`
		MaxActiveCollaboratorPerMonth int       `db:"max_active_collaborator_per_month"`
		UpdatedAt                     time.Time `db:"updated_at"`
	}{
		ID:                            p.ID,
		Name:                          p.Name,
		SendMonthlyReport:             p.SendMonthlyReport,
		CanCustomizeReport:            p.CanCustomizeReport,
		SendExpenseReport:             p.SendExpenseReport,
		MaxActiveCollaboratorPerMonth: p.MaxActiveCollaboratorPerMonth,
		UpdatedAt:                     now,
	}

	const q = `
		UPDATE 
			"public"."enterprise_pack" 
		SET 
			name 								= :name 
			send_expense_report 				= :send_expense_report
			can_customize_report 				= :can_customize_report
			send_expense_report 				= :send_expense_report
			max_active_collaborator_per_month 	= :max_active_collaborator_per_month
			updated_at         	 				= :updated_at
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
		"public"."enterprise_pack"
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
