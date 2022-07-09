package paymentmethod

import (
	"context"
	"database/sql"
	"fmt"
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

func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

func (s Store) Create(ctx context.Context, stripeID string, npm dto.NewPaymentMethod, now time.Time) (PaymentMethod, error) {
	p := PaymentMethod{
		ID:                validate.GenerateID(),
		Name:              npm.Name,
		UserID:            npm.UserID,
		DisplayCreditCard: fmt.Sprintf("XXXX XXXX XXXX %s", npm.Number[len(npm.Number)-4:]),
		Type:              npm.Type,
		IsFavorite:        npm.IsFavorite,
		StripeID:          stripeID,
		UpdatedAt:         now,
		CreatedAt:         now,
		DeletedAt:         pq.NullTime{},
	}

	const q = `
	INSERT INTO "public"."payment_method"
		(id, user_id, stripe_id, name, display_credit_card, type, is_favorite, created_at, updated_at, deleted_at)
	VALUES 
		(:id, :user_id, :stripe_id, :name, :display_credit_card, :type, :is_favorite, :created_at, :updated_at, :deleted_at)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, p); err != nil {
		return PaymentMethod{}, err
	}

	return p, nil
}

//Query get all the payment methods of a specified user. For restriction reason
//we allow to query payment methods only per user
func (s Store) Query(ctx context.Context, id string, pageNumber, rowsPerPage int) ([]PaymentMethod, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		UserID      string `db:"user_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		UserID:      id,
	}

	const q = `
		SELECT 
			id, 
			name, 
			user_id, 
			display_credit_card, 
			stripe_id, 
			type, 
			is_favorite, 
			created_at, 
			updated_at
		FROM "public"."payment_method"
			WHERE deleted_at IS NULL 
			AND user_id = :user_id
		ORDER BY 
			id
		OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	pms := make([]PaymentMethod, rowsPerPage)

	if err := database.NamedQuerySlice[PaymentMethod](ctx, s.log, s.db, q, data, &pms); err != nil {
		return pms, err
	}

	return pms, nil
}

func (s Store) QueryByID(ctx context.Context, id string) (PaymentMethod, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
		SELECT 
			id, 
			name, 
			user_id, 
			display_credit_card, 
			stripe_id, 
			type, 
			is_favorite, 
			created_at, 
			updated_at
		FROM "public"."payment_method"
			WHERE deleted_at IS NULL 
			AND id = :id`

	var p PaymentMethod

	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &p); err != nil {
		return PaymentMethod{}, err
	}

	return p, nil
}

func (s Store) Update(ctx context.Context, p PaymentMethod, now time.Time) error {
	data := struct {
		ID         string    `db:"id"`
		IsFavorite bool      `db:"is_favorite"`
		Name       string    `db:"name"`
		UpdatedAt  time.Time `db:"updated_at"`
	}{
		ID:         p.ID,
		IsFavorite: p.IsFavorite,
		Name:       p.Name,
		UpdatedAt:  now,
	}

	const q = `
	UPDATE 
		"public"."payment_method"
	SET 
		is_favorite = :is_favorite
		name 		= :name
		updated_at  = :updated_at
	WHERE 
		id 			= :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}

func (s Store) Delete(ctx context.Context, id string, now time.Time) error {
	data := struct {
		ID        string       `db:"id"`
		DeletedAt sql.NullTime `db:"deleted_at"`
	}{
		ID: id,
		DeletedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	const q = `
	UPDATE 
		"public"."payment_method"
	SET
		deleted_at  = :deleted_at
	WHERE 
		id 			= :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return err
	}

	return nil
}
