package user

import (
	"context"
	"fmt"
	"time"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/internal/repo"
	"github.com/71g3pf4c3/gophermart/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
	jwt  *jwt.Manager
}

// New -.
func New(r repo.UserRepo, j *jwt.Manager) *UseCase {
	return &UseCase{
		repo: r,
		jwt:  j,
	}
}

// Register -.
func (uc *UseCase) Register(ctx context.Context, username, password string) (entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	now := time.Now().UTC()

	user := entity.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    now,
	}

	err = uc.repo.CreateUser(ctx, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.CreateUser: %w", err)
	}

	return user, nil
}

// Login -.
func (uc *UseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Login - uc.jwt.GenerateToken: %w", err)
	}

	return token, nil
}

// GetUser -.
func (uc *UseCase) GetUser(ctx context.Context, userID string) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetByID: %w", err)
	}

	return user, nil
}
