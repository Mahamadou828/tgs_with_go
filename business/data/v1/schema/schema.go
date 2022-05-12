//Package schema allow database migration, seeding and
//contains schema
package schema

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"

	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/v1/schema.sql
	schemaV1Doc string
	//go:embed sql/v1/seed.sql
	seedV1Doc string
)

type Migration struct {
	Schema string
	Seed   string
}

var migration = make(map[string]Migration)

func init() {
	migration["v1"] = Migration{
		Schema: schemaV1Doc,
		Seed:   seedV1Doc,
	}
}

//Migrate attempts to bring the schema for db up to date with the migration
//defined in this package
func Migrate(ctx context.Context, db *sqlx.DB, version string) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("database status check failed: %v", err)
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	query, ok := migration[version]

	if !ok {
		return fmt.Errorf("unavailable schema query for version %s", version)
	}

	if _, err := tx.Exec(query.Schema); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func Seed(ctx context.Context, db *sqlx.DB, version string) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("database status check failed: %v", err)
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	query, ok := migration[version]

	if !ok {
		return fmt.Errorf("unavailable seed query for version %s", version)
	}

	if _, err := tx.Exec(query.Seed); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
