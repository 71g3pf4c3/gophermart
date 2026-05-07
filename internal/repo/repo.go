package repo

import (
	"context"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

type UserRepo interface {
		CreateUser(ctx context.Context, user *entity.User) (error)
		GetByID(ctx context.Context, userID string) (entity.User, error)
		GetByUsername(ctx context.Context, username string) (entity.User, error)
	}
