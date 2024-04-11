package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-logr/logr"
	_ "github.com/lib/pq"
)

type Database struct {
	conn   *sql.DB
	logger logr.Logger
}

func (db *Database) String() string {
	return "postgres-connection"
}

func (db *Database) GracefulShutdown(_ context.Context) error {
	return db.conn.Close() //nolint:wrapcheck
}

func (db *Database) WithinTransaction(ctx context.Context, wrappedFunc func(ctx context.Context) error) error {
	tran, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction:%w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			db.rollbackTransaction(tran)
			panic(r)
		}
	}()

	if err = wrappedFunc(saveTxToContext(ctx, tran)); err != nil {
		db.rollbackTransaction(tran)
		return fmt.Errorf("within transaction: %w", err)
	}

	if err := tran.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (db *Database) rollbackTransaction(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		db.logger.Error(err, "failed to rollback transaction")
	}
}

//nolint:wrapcheck
func (db *Database) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx.PrepareContext(ctx, query)
	}
	return db.conn.PrepareContext(ctx, query)
}

//nolint:wrapcheck
func (db *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return db.conn.ExecContext(ctx, query, args...)
}

//nolint:wrapcheck
func (db *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return db.conn.QueryRowContext(ctx, query, args...)
}

func (db *Database) Postgres() *sql.DB {
	return db.conn
}

func New(logger logr.Logger, url string) (*Database, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Database{
		conn:   conn,
		logger: logger,
	}, nil
}

func NewFromConnection(logger logr.Logger, conn *sql.DB) *Database {
	return &Database{
		conn:   conn,
		logger: logger,
	}
}
