package balance

import (
	"context"
	"errors"
	"testing"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

type fakeBalanceRepo struct {
	balance entity.Balance
	w       []entity.Withdrawal
	err     error
}

func (r *fakeBalanceRepo) GetBalance(_ context.Context, _ string) (entity.Balance, error) {
	return r.balance, nil
}

func (r *fakeBalanceRepo) Withdraw(_ context.Context, _, _ string, _ float64) error {
	return r.err
}

func (r *fakeBalanceRepo) GetWithdrawals(_ context.Context, _ string) ([]entity.Withdrawal, error) {
	return r.w, nil
}

func TestBalanceUseCase(t *testing.T) {
	repo := &fakeBalanceRepo{balance: entity.Balance{Current: 10, Withdrawn: 5}}
	uc := New(repo)

	b, err := uc.GetBalance(context.Background(), "u1")
	if err != nil || b.Current != 10 {
		t.Fatalf("unexpected balance result: %+v err=%v", b, err)
	}

	repo.err = entity.ErrInsufficientFunds
	err = uc.Withdraw(context.Background(), "u1", "2377225624", 100)
	if !errors.Is(err, entity.ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}

	err = uc.Withdraw(context.Background(), "u1", "79927398710", 1)
	if !errors.Is(err, entity.ErrInvalidOrderNumber) {
		t.Fatalf("expected ErrInvalidOrderNumber, got %v", err)
	}
}
