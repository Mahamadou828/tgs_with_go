package database

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

//Database define an instance of the database connection
type Database struct {
	*sqlx.DB
}

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

//Open create a new database connection
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"

	if cfg.DisableTLS {
		sslMode = "disable"
	}
	q := make(url.Values)

	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())

	if err != nil {
		return nil, fmt.Errorf("can't connect to database received: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database record: %v", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	err := ctx.Err()

	// Make sure we didn't timeout or be cancelled.
	if err != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// NamedExecContext is a helper function to execute a CRUD operation with
// logging and tracing.
func NamedExecContext(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data interface{}) error {
	q := queryString(query, ctx)
	log.Infow("database.NamedExecContext", "traceId", web.GetTraceID(ctx), "query", q)

	if _, err := sqlx.NamedExecContext(ctx, db, q, data); err != nil {
		return err
	}
	return nil
}

// NamedQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshalled into a slice.
func NamedQuerySlice[T any](ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data interface{}, dest *[]T) error {
	q := queryString(query, data)

	log.Infow("database.NamedQuerySlice", "traceId", web.GetTraceID(ctx), "query", q)

	rows, err := sqlx.NamedQueryContext(ctx, db, q, db)

	if err != nil {
		return err
	}

	defer rows.Close()

	var slice []T

	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}
	*dest = slice

	return nil
}

func queryString(query string, args ...any) string {
	query, params, err := sqlx.Named(query, args)

	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string

		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\t", " ")

	return query
}
