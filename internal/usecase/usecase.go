package usecase

import (
	"context"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

type (
	User interface {
		Register(ctx context.Context, login, password string) (entity.User, string, error)
		Login(ctx context.Context, login, password string) (string, error)
		GetUser(ctx context.Context, userID string) (entity.User, error)
	}
	TokenManager interface {
		GenerateToken(userID string) (string, error)
		ParseToken(token string) (string, error)
	}
)
