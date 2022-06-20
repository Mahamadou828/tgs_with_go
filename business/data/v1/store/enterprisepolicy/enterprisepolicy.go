package enterprisepolicy

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
	db  *sqlx.DB
	log *zap.SugaredLogger
}

func NewStore(db *sqlx.DB, log *zap.SugaredLogger) Store {
	return Store{
		db:  db,
		log: log,
	}
}

func (s Store) Create(ctx context.Context, ntp dto.NewEnterprisePolicy, now time.Time) (TeamPolicy, error) {
	tp := TeamPolicy{
		ID:                 validate.GenerateID(),
		Name:               ntp.Name,
		Description:        ntp.Description,
		CollaboratorBudget: ntp.CollaboratorBudget,
		StartServiceTime:   ntp.StartServiceTime,
		EndServiceTime:     ntp.EndServiceTime,
		BudgetType:         ntp.BudgetType,
		EnterpriseID:       ntp.EnterpriseID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	const q = `
	INSERT INTO enterprise_policy 
		(id, enterprise_id, name, blocked_days, description, collaborator_budget, start_service_time, end_service_time, budget_type, created_at, updated_at, deleted_at)
	VALUES
		(:id, :enterprise_id, :name, :blocked_days, :description, :collaborator_budget, :start_service_time, :end_service_time, :budget_type, :created_at, :updated_at, null)
`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, tp); err != nil {
		return TeamPolicy{}, err
	}
	return tp, nil
}

func (s Store) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]TeamPolicy, error) {
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
		description, 
		collaborator_budget, 
		start_service_time, 
		end_service_time, 
		budget_type, 
		enterprise_id,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise_policy" 
	WHERE deleted_at IS NULL
	ORDER BY 
		id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	policy := make([]TeamPolicy, rowsPerPage)

	if err := database.NamedQuerySlice[TeamPolicy](ctx, s.log, s.db, q, data, &policy); err != nil {
		return []TeamPolicy{}, err
	}

	return policy, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (TeamPolicy, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		id, 
		name,
		description, 
		collaborator_budget, 
		start_service_time, 
		end_service_time, 
		budget_type, 
		enterprise_id,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise_policy" 
	WHERE deleted_at IS NULL AND id = :id 
`

	var team TeamPolicy

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return TeamPolicy{}, err
	}

	return team, nil
}

func (s Store) QueryByEnterpriseID(ctx context.Context, id string) (TeamPolicy, error) {
	data := struct {
		EnterpriseID string `db:"enterprise_id"`
	}{
		EnterpriseID: id,
	}

	const q = `
	SELECT 
		id, 
		name,
		description, 
		collaborator_budget, 
		start_service_time, 
		end_service_time, 
		budget_type, 
		enterprise_id,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise_policy" 
	WHERE deleted_at IS NULL AND enterprise_id = :enterprise_id 
`

	var team TeamPolicy

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return TeamPolicy{}, err
	}

	return team, nil
}

func (s Store) QueryByTeamID(ctx context.Context, id string) (TeamPolicy, error) {
	data := struct {
		TeamID string `db:"team_id"`
	}{
		TeamID: id,
	}

	const q = `
	SELECT 
		id, 
		name,
		description, 
		collaborator_budget, 
		start_service_time, 
		end_service_time, 
		budget_type, 
		enterprise_id,
		updated_at, 
		created_at, 
		deleted_at
	FROM 
		"public"."enterprise_policy" 
	WHERE deleted_at IS NULL AND team_id = :team_id 
`

	var team TeamPolicy

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return TeamPolicy{}, err
	}

	return team, nil
}

func (s Store) Update(ctx context.Context, id string, tp TeamPolicy, now time.Time) error {
	data := struct {
		UpdatedAt          time.Time `db:"updated_at"`
		ID                 string    `db:"id"`
		Name               string    `db:"name"`
		Description        string    `db:"description"`
		CollaboratorBudget int       `db:"collaborator_budget"`
		StartServiceTime   string    `db:"start_service_time"`
		EndServiceTime     string    `db:"end_service_time"`
		BudgetType         string    `db:"budget_type"`
	}{
		UpdatedAt:          now,
		ID:                 id,
		Name:               tp.Name,
		Description:        tp.Description,
		CollaboratorBudget: tp.CollaboratorBudget,
		StartServiceTime:   tp.StartServiceTime,
		EndServiceTime:     tp.EndServiceTime,
		BudgetType:         tp.BudgetType,
	}

	const q = `
	UPDATE 
		"public"."enterprise_policy" 
	SET
		description         = :description,
		name 				= :name,
		collaborator_budget = :collaborator_budget,
		start_service_time 	= :start_service_time,
		end_service_time 	= :end_service_time,
		budget_type 		= :budget_type
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
		"public"."enterprise_policy"
	SET
		deleted_at 		= :deleted_at
	WHERE 	
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
