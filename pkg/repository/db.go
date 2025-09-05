package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

var (
	ErrTransactionClosed = errors.New("sql: transaction has already been committed or rolled back")
)

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error)
	Dialect() schema.Dialect
	NewValues(model interface{}) *bun.ValuesQuery
	NewMerge() *bun.MergeQuery
	NewSelect() *bun.SelectQuery
	NewInsert() *bun.InsertQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
	NewRaw(query string, args ...interface{}) *bun.RawQuery
	NewCreateTable() *bun.CreateTableQuery
	NewDropTable() *bun.DropTableQuery
	NewCreateIndex() *bun.CreateIndexQuery
	NewDropIndex() *bun.DropIndexQuery
	NewTruncateTable() *bun.TruncateTableQuery
	NewAddColumn() *bun.AddColumnQuery
	NewDropColumn() *bun.DropColumnQuery
}

// type TransactionCloser interface {
// 	Commit() error
// 	Rollback() error
// }

// type InTransactionFunc func(ctx context.Context, db DB) (bool, error)

// type errorRow struct {
// 	err error
// }

// func (r *errorRow) Scan(dest ...any) error {
// 	return r.err
// }

// func (r *errorRow) Err() error {
// 	return r.err
// }

// func InTransaction(ctx context.Context, db DB, fn InTransactionFunc) error {
// 	db, tx, err := db.BeginTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	shouldCommit, err := fn(ctx, db)
// 	if err != nil {
// 		return err
// 	}
// 	if shouldCommit {
// 		err = tx.Commit()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func NewDB(db *sql.DB) DB {
// 	return &sqlDbWrapper{db: db}
// }

// type sqlDbWrapper struct {
// 	db *sql.DB
// }

// func (d *sqlDbWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
// 	return d.db.ExecContext(ctx, query, args...)
// }

// func (d *sqlDbWrapper) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
// 	return d.db.PrepareContext(ctx, query)
// }

// func (d *sqlDbWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
// 	return d.db.QueryContext(ctx, query, args...)
// }

// func (d *sqlDbWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
// 	return d.db.QueryRowContext(ctx, query, args...)
// }

// func (d *sqlDbWrapper) BeginTx(ctx context.Context) (DB, TransactionCloser, error) {
// 	tx, err := d.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return d, nil, err // important to return the actual DB in case usage is `db, err = db.BeginTx(ctx)`
// 	}

// 	newDb := &dbInTransaction{tx: tx}

// 	return newDb, newDb, nil
// }

// type dbInTransaction struct {
// 	tx *sql.Tx
// }

// func (d *dbInTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
// 	return d.tx.ExecContext(ctx, query, args...)
// }

// func (d *dbInTransaction) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
// 	return d.tx.PrepareContext(ctx, query)
// }

// func (d *dbInTransaction) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
// 	return d.tx.QueryContext(ctx, query, args...)
// }

// func (d *dbInTransaction) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
// 	return d.tx.QueryRowContext(ctx, query, args...)
// }

// func (d *dbInTransaction) BeginTx(ctx context.Context) (DB, TransactionCloser, error) {
// 	return createSavepointDB(ctx, d)
// }

// func (d *dbInTransaction) Rollback() error {
// 	return d.tx.Rollback()
// }

// func (d *dbInTransaction) Commit() error {
// 	return d.tx.Commit()
// }

// func createSavepointDB(ctx context.Context, baseDb DB) (DB, TransactionCloser, error) {
// 	savepointUuid, err := uuid.NewV7()
// 	if err != nil {
// 		// important to return the actual DB in case usage is `db, err = db.BeginTx(ctx)`
// 		return baseDb, nil, fmt.Errorf("creating a new savepoint unique ID: %w", err)
// 	}

// 	newDb := &dbInSavepoint{
// 		ctx:           ctx,
// 		baseDb:        baseDb,
// 		savepointName: "sp_" + savepointUuid.String(),
// 	}

// 	_, err = baseDb.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %q", newDb.savepointName))
// 	if err != nil {
// 		// important to return the actual DB in case usage is `db, err = db.BeginTx(ctx)`
// 		return baseDb, nil, fmt.Errorf("executing savepoint creation: %w", err)
// 	}

// 	return newDb, newDb, nil
// }

// type dbInSavepoint struct {
// 	baseDb        DB
// 	savepointName string
// 	ctx           context.Context
// 	lock          sync.RWMutex
// 	closed        bool
// }

// func (d *dbInSavepoint) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	if d.closed {
// 		return nil, ErrTransactionClosed
// 	}

// 	return d.baseDb.ExecContext(ctx, query, args...)
// }

// func (d *dbInSavepoint) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	if d.closed {
// 		return nil, ErrTransactionClosed
// 	}

// 	return d.baseDb.PrepareContext(ctx, query)
// }

// func (d *dbInSavepoint) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	if d.closed {
// 		return nil, ErrTransactionClosed
// 	}

// 	return d.baseDb.QueryContext(ctx, query, args...)
// }

// func (d *dbInSavepoint) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	if d.closed {
// 		panic(ErrTransactionClosed)
// 	}

// 	return d.baseDb.QueryRowContext(ctx, query, args...)
// }

// func (d *dbInSavepoint) BeginTx(ctx context.Context) (DB, TransactionCloser, error) {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	if d.closed {
// 		return nil, nil, ErrTransactionClosed
// 	}

// 	return createSavepointDB(ctx, d)
// }

// func (d *dbInSavepoint) Rollback() error {
// 	d.lock.Lock()
// 	defer d.lock.Unlock()

// 	if d.closed {
// 		return nil
// 	}

// 	_, err := d.baseDb.ExecContext(d.ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %q", d.savepointName))
// 	if err != nil {
// 		return fmt.Errorf("executing savepoint %q release (commit): %w", d.savepointName, err)
// 	}

// 	d.closed = true
// 	return nil
// }

// func (d *dbInSavepoint) Commit() error {
// 	d.lock.Lock()
// 	defer d.lock.Unlock()

// 	if d.closed {
// 		return nil
// 	}

// 	_, err := d.baseDb.ExecContext(d.ctx, fmt.Sprintf("RELEASE SAVEPOINT %q", d.savepointName))
// 	if err != nil {
// 		return fmt.Errorf("executing savepoint %q rollback: %w", d.savepointName, err)
// 	}

// 	d.closed = true
// 	return nil
// }
