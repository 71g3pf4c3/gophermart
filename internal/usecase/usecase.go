package usecase

import (
	"context"

	"github.com/71g3pf4c3/gophermart/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// User -.
	User interface {
		Register(ctx context.Context, username, password string) (entity.User, error)
		Login(ctx context.Context, username, password string) (string, error)
		GetUser(ctx context.Context, userID string) (entity.User, error)
	}
	TokenManager interface {
		GenerateToken(userID string) (string, error)
		ParseToken(token string) (string, error)
	}
)
