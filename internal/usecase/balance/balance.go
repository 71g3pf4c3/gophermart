package balance

import (
	"context"
	"fmt"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/luhn"
)

// BalanceRepo is the required repository interface.
type BalanceRepo interface {
	GetBalance(ctx context.Context, userID string) (entity.Balance, error)
	Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}

// UseCase handles balance and withdrawal business logic.
type UseCase struct {
	repo BalanceRepo
}

// New creates a new balance UseCase.
func New(r BalanceRepo) *UseCase {
	return &UseCase{repo: r}
}

// GetBalance returns current balance and total withdrawn.
func (uc *UseCase) GetBalance(ctx context.Context, userID string) (entity.Balance, error) {
	b, err := uc.repo.GetBalance(ctx, userID)
	if err != nil {
		return entity.Balance{}, fmt.Errorf("balance - GetBalance: %w", err)
	}
	return b, nil
}

// Withdraw deducts sum from user balance for the given order number.
func (uc *UseCase) Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error {
	if !luhn.Valid(orderNumber) {
		return entity.ErrInvalidOrderNumber
	}
	return uc.repo.Withdraw(ctx, userID, orderNumber, sum)
}

// GetWithdrawals returns all withdrawals for a user.
func (uc *UseCase) GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdrawal, error) {
	ws, err := uc.repo.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("balance - GetWithdrawals: %w", err)
	}
	return ws, nil
}
