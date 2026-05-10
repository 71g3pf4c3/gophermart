package repo

import (
	"context"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, userID string) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}
