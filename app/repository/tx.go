package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager struct {
	pool      *pgxpool.Pool
	runTxOnce func(ctx context.Context, fn func(tx DBTX) error) error
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	m := &TxManager{pool: pool}
	m.runTxOnce = m.runTransactionOnce
	return m
}

// WithTx menjalankan callback dalam transaksi dengan retry hingga 3 kali
// (backoff 50/100/200 ms) HANYA untuk error PG 40001 (serialization) dan 40P01 (deadlock).
func (m *TxManager) WithTx(ctx context.Context, fn func(tx DBTX) error) error {
	backoffs := []time.Duration{50 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond}
	maxAttempts := len(backoffs) + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := m.runTxOnce(ctx, fn)
		if err == nil {
			return nil
		}

		if isRetryablePgError(err) {
			if attempt < maxAttempts {
				select {
				case <-time.After(backoffs[attempt-1]):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		return err
	}
	return nil
}

// WithTxOnce menjalankan callback transaksi tepat satu kali tanpa retry.
func (m *TxManager) WithTxOnce(ctx context.Context, fn func(tx DBTX) error) error {
	return m.runTxOnce(ctx, fn)
}

func isRetryablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001" || pgErr.Code == "40P01"
	}
	return false
}

func (m *TxManager) runTransactionOnce(ctx context.Context, fn func(tx DBTX) error) (err error) {
	if m.pool == nil {
		return errors.New("database pool is not initialized")
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			rollbackErr = fmt.Errorf("rollback transaction: %w", rollbackErr)
			if err == nil {
				err = rollbackErr
			} else {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	if err := fn(tx); err != nil {
		return fmt.Errorf("transaction callback: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
