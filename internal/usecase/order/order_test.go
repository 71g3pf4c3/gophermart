package order

import (
	"context"
	"errors"
	"testing"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

type fakeOrderRepo struct{ orders map[string]entity.Order }

func (r *fakeOrderRepo) CreateOrder(_ context.Context, order entity.Order) error {
	r.orders[order.Number] = order
	return nil
}

func (r *fakeOrderRepo) GetOrderByNumber(_ context.Context, number string) (entity.Order, error) {
	o, ok := r.orders[number]
	if !ok {
		return entity.Order{}, entity.ErrOrderNotFound
	}
	return o, nil
}

func (r *fakeOrderRepo) GetOrdersByUser(_ context.Context, userID string) ([]entity.Order, error) {
	var out []entity.Order
	for _, o := range r.orders {
		if o.UserID == userID {
			out = append(out, o)
		}
	}
	return out, nil
}

func TestUploadAndList(t *testing.T) {
	repo := &fakeOrderRepo{orders: map[string]entity.Order{}}
	uc := New(repo)

	if err := uc.Upload(context.Background(), "u1", "79927398713"); err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if err := uc.Upload(context.Background(), "u1", "79927398713"); !errors.Is(err, entity.ErrOrderOwnedBySelf) {
		t.Fatalf("expected ErrOrderOwnedBySelf, got %v", err)
	}
	if err := uc.Upload(context.Background(), "u2", "79927398713"); !errors.Is(err, entity.ErrOrderOwnedByOther) {
		t.Fatalf("expected ErrOrderOwnedByOther, got %v", err)
	}
	if err := uc.Upload(context.Background(), "u1", "79927398710"); !errors.Is(err, entity.ErrInvalidOrderNumber) {
		t.Fatalf("expected ErrInvalidOrderNumber, got %v", err)
	}

	orders, err := uc.List(context.Background(), "u1")
	if err != nil || len(orders) != 1 {
		t.Fatalf("expected one order, got len=%d err=%v", len(orders), err)
	}
}
