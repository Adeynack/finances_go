package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Transaction(ctx context.Context, fn InTransactionFn) error
}

type InTransactionFn func(ctx context.Context, db DB) (bool, error)

func NewDB(db *sql.DB) DB {
	return &sqlDbWrapper{db: db}
}

type sqlDbWrapper struct {
	db *sql.DB
}

func (d *sqlDbWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *sqlDbWrapper) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.db.PrepareContext(ctx, query)
}

func (d *sqlDbWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *sqlDbWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *sqlDbWrapper) Transaction(ctx context.Context, fn InTransactionFn) (err error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}

	defer func() {
		defErr := tx.Rollback()
		err = errors.Join(err, defErr)
	}()

	newDb := &dbInTransaction{tx: tx}

	shouldCommit, err := fn(ctx, newDb)
	if err != nil {
		return err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

type dbInTransaction struct {
	tx *sql.Tx
}

func (d *dbInTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.tx.ExecContext(ctx, query, args...)
}

func (d *dbInTransaction) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.tx.PrepareContext(ctx, query)
}

func (d *dbInTransaction) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.tx.QueryContext(ctx, query, args...)
}

func (d *dbInTransaction) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.tx.QueryRowContext(ctx, query, args...)
}

func (d *dbInTransaction) Transaction(ctx context.Context, fn InTransactionFn) error {
	savepointUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	savepointName := "sp_" + savepointUuid.String()

	_, err = d.tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %q", savepointName))
	if err != nil {
		return err
	}
	shouldCommit := false

	defer func() {
		var derEff error
		if shouldCommit {
			_, derEff = d.tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %q", savepointName))
		} else {
			_, derEff = d.tx.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %q", savepointName))
		}
		err = errors.Join(err, derEff)
	}()

	shouldCommit, err = fn(ctx, d)
	if err != nil {
		return err
	}

	return nil
}
