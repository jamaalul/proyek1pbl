package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestTxManager_RetryLogic(t *testing.T) {
	t.Run("SuccessOnFirstAttempt", func(t *testing.T) {
		m := NewTxManager(nil)
		attempts := 0
		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			return nil
		}

		err := m.WithTx(context.Background(), func(tx DBTX) error {
			return nil
		})

		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got: %d", attempts)
		}
	})

	t.Run("NonRetryableErrorFailsImmediately", func(t *testing.T) {
		m := NewTxManager(nil)
		attempts := 0
		nonRetryErr := &pgconn.PgError{Code: "23505", Message: "duplicate key"}
		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			return nonRetryErr
		}

		err := m.WithTx(context.Background(), func(tx DBTX) error {
			return nil
		})

		if !errors.Is(err, nonRetryErr) {
			t.Errorf("expected nonRetryErr, got: %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected exactly 1 attempt without retry, got: %d", attempts)
		}
	})

	t.Run("RetryOnSerializationFailure_SucceedsEventually", func(t *testing.T) {
		m := NewTxManager(nil)
		attempts := 0
		serializationErr := &pgconn.PgError{Code: "40001", Message: "serialization_failure"}

		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			if attempts < 3 {
				return serializationErr
			}
			return nil
		}

		err := m.WithTx(context.Background(), func(tx DBTX) error {
			return nil
		})

		if err != nil {
			t.Fatalf("expected eventual success, got: %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got: %d", attempts)
		}
	})

	t.Run("RetryOnDeadlock_GivesUpAfterMaxAttempts", func(t *testing.T) {
		m := NewTxManager(nil)
		attempts := 0
		deadlockErr := &pgconn.PgError{Code: "40P01", Message: "deadlock_detected"}

		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			return deadlockErr
		}

		err := m.WithTx(context.Background(), func(tx DBTX) error {
			return nil
		})

		if !errors.Is(err, deadlockErr) {
			t.Errorf("expected deadlockErr after exhaustion, got: %v", err)
		}
		// 1 initial attempt + 3 retries = 4 attempts total
		if attempts != 4 {
			t.Errorf("expected 4 attempts total (1 initial + 3 retries), got: %d", attempts)
		}
	})

	t.Run("WithTxOnceDoesNotRetry", func(t *testing.T) {
		m := NewTxManager(nil)
		attempts := 0
		deadlockErr := &pgconn.PgError{Code: "40P01", Message: "deadlock_detected"}

		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			return deadlockErr
		}

		err := m.WithTxOnce(context.Background(), func(tx DBTX) error {
			return nil
		})

		if !errors.Is(err, deadlockErr) {
			t.Errorf("expected deadlockErr, got: %v", err)
		}
		if attempts != 1 {
			t.Errorf("WithTxOnce should only execute 1 attempt, got: %d", attempts)
		}
	})

	t.Run("ContextCancellationAbortsRetry", func(t *testing.T) {
		m := NewTxManager(nil)
		ctx, cancel := context.WithCancel(context.Background())
		attempts := 0
		serializationErr := &pgconn.PgError{Code: "40001"}

		m.runTxOnce = func(ctx context.Context, fn func(tx DBTX) error) error {
			attempts++
			cancel() // cancel context on first failure
			return serializationErr
		}

		err := m.WithTx(ctx, func(tx DBTX) error {
			return nil
		})

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got: %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt before cancellation, got: %d", attempts)
		}
	})

	t.Run("NilPoolReturnsError", func(t *testing.T) {
		m := NewTxManager(nil)
		// test default runTransactionOnce with nil pool
		err := m.runTransactionOnce(context.Background(), func(tx DBTX) error {
			return nil
		})
		if err == nil || err.Error() != "database pool is not initialized" {
			t.Errorf("expected 'database pool is not initialized', got: %v", err)
		}
	})
}
