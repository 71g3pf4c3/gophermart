package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/jwt"
)

type fakeUserRepo struct{ users map[string]entity.User }

func (r *fakeUserRepo) CreateUser(_ context.Context, user *entity.User) error {
	if _, ok := r.users[user.Login]; ok {
		return entity.ErrUserAlreadyExists
	}
	r.users[user.Login] = *user
	return nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, userID string) (entity.User, error) {
	for _, u := range r.users {
		if u.ID == userID {
			return u, nil
		}
	}
	return entity.User{}, entity.ErrUserNotFound
}

func (r *fakeUserRepo) GetByLogin(_ context.Context, login string) (entity.User, error) {
	u, ok := r.users[login]
	if !ok {
		return entity.User{}, entity.ErrUserNotFound
	}
	return u, nil
}

func TestRegisterAndLogin(t *testing.T) {
	repo := &fakeUserRepo{users: map[string]entity.User{}}
	uc := New(repo, jwt.New("secret", time.Hour))

	user, token, err := uc.Register(context.Background(), "alice", "pass123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if user.Login != "alice" || token == "" {
		t.Fatal("unexpected register result")
	}

	_, _, err = uc.Register(context.Background(), "alice", "pass123")
	if !errors.Is(err, entity.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}

	_, err = uc.Login(context.Background(), "alice", "badpass")
	if !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	goodToken, err := uc.Login(context.Background(), "alice", "pass123")
	if err != nil || goodToken == "" {
		t.Fatalf("expected successful login, err=%v token=%q", err, goodToken)
	}
}
