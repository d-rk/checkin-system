package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tx is an interface that models the standard transaction in `database/sql`.
//
// To ensure `TransactionalFunc` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Tx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

// TransactionalFunc a function that will be called with an initialized `Tx` object
// that can be used for executing statements and queries against a database.
type TransactionalFunc func(Tx) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TransactionalFunc`.
func WithTransaction(db *sqlx.DB, fn TransactionalFunc) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
