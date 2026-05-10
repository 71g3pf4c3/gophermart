package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// BalanceRepository implements balance and withdrawal storage.
type BalanceRepository struct {
	*postgres.Postgres
}

// NewBalanceRepository creates a new BalanceRepository.
func NewBalanceRepository(pg *postgres.Postgres) *BalanceRepository {
	return &BalanceRepository{pg}
}

// GetBalance returns current balance and total withdrawn for a user.
func (r *BalanceRepository) GetBalance(ctx context.Context, userID string) (entity.Balance, error) {
	const q = `
		SELECT
			COALESCE((SELECT SUM(accrual) FROM orders WHERE user_id = $1 AND status = 'PROCESSED'), 0)
				- COALESCE((SELECT SUM(sum) FROM withdrawals WHERE user_id = $1), 0),
			COALESCE((SELECT SUM(sum) FROM withdrawals WHERE user_id = $1), 0)
	`

	var b entity.Balance
	err := r.Pool.QueryRow(ctx, q, userID).Scan(&b.Current, &b.Withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Balance{}, nil
		}
		return entity.Balance{}, fmt.Errorf("BalanceRepository - GetBalance - Scan: %w", err)
	}

	return b, nil
}

// Withdraw creates a withdrawal record inside a transaction with an advisory lock.
func (r *BalanceRepository) Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("BalanceRepository - Withdraw - Begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Advisory lock per user prevents concurrent overdraw
	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, userID); err != nil {
		return fmt.Errorf("BalanceRepository - Withdraw - advisory lock: %w", err)
	}

	const balanceQ = `
		SELECT
			COALESCE((SELECT SUM(accrual) FROM orders WHERE user_id = $1 AND status = 'PROCESSED'), 0)
				- COALESCE((SELECT SUM(sum) FROM withdrawals WHERE user_id = $1), 0)
	`

	var current float64
	if err = tx.QueryRow(ctx, balanceQ, userID).Scan(&current); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("BalanceRepository - Withdraw - balance query: %w", err)
	}

	if current < sum {
		return entity.ErrInsufficientFunds
	}

	const insertQ = `INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`
	if _, err = tx.Exec(ctx, insertQ, userID, orderNumber, sum); err != nil {
		return fmt.Errorf("BalanceRepository - Withdraw - insert: %w", err)
	}

	return tx.Commit(ctx)
}

// GetWithdrawals returns withdrawals for a user sorted by processed_at desc.
func (r *BalanceRepository) GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdrawal, error) {
	const q = `
		SELECT order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`

	rows, err := r.Pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("BalanceRepository - GetWithdrawals - Query: %w", err)
	}
	defer rows.Close()

	var withdrawals []entity.Withdrawal
	for rows.Next() {
		var w entity.Withdrawal
		if err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, fmt.Errorf("BalanceRepository - GetWithdrawals - Scan: %w", err)
		}
		withdrawals = append(withdrawals, w)
	}

	return withdrawals, rows.Err()
}
