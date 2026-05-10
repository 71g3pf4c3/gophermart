package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/luhn"
)

// OrderRepo is the required repository interface.
type OrderRepo interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	GetOrderByNumber(ctx context.Context, number string) (entity.Order, error)
	GetOrdersByUser(ctx context.Context, userID string) ([]entity.Order, error)
}

// UseCase handles order business logic.
type UseCase struct {
	repo OrderRepo
}

// New creates a new order UseCase.
func New(r OrderRepo) *UseCase {
	return &UseCase{repo: r}
}

// Upload accepts a new order number for processing.
func (uc *UseCase) Upload(ctx context.Context, userID, number string) error {
	if !luhn.Valid(number) {
		return entity.ErrInvalidOrderNumber
	}

	existing, err := uc.repo.GetOrderByNumber(ctx, number)
	if err != nil && !errors.Is(err, entity.ErrOrderNotFound) {
		return fmt.Errorf("order - Upload - GetOrderByNumber: %w", err)
	}

	if err == nil {
		if existing.UserID == userID {
			return entity.ErrOrderOwnedBySelf
		}
		return entity.ErrOrderOwnedByOther
	}

	return uc.repo.CreateOrder(ctx, entity.Order{
		Number:     number,
		UserID:     userID,
		Status:     entity.OrderStatusNew,
		UploadedAt: time.Now().UTC(),
	})
}

// List returns all orders for a user.
func (uc *UseCase) List(ctx context.Context, userID string) ([]entity.Order, error) {
	orders, err := uc.repo.GetOrdersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("order - List - GetOrdersByUser: %w", err)
	}
	return orders, nil
}
