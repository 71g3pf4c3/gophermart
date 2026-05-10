package restapi

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type fakeUserService struct {
	regErr error
	logErr error
}

func (f *fakeUserService) Register(_ context.Context, _, _ string) (entity.User, string, error) {
	if f.regErr != nil {
		return entity.User{}, "", f.regErr
	}
	return entity.User{ID: "u1", Login: "alice"}, "tok", nil
}

func (f *fakeUserService) Login(_ context.Context, _, _ string) (string, error) {
	if f.logErr != nil {
		return "", f.logErr
	}
	return "tok", nil
}

func TestUserHandlers(t *testing.T) {
	app := fiber.New()
	h := newUserHandler(&fakeUserService{})
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)

	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"login":"alice","password":"pass"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("register status=%d", resp.StatusCode)
	}

	req = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"login":"alice","password":"pass"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("login status=%d", resp.StatusCode)
	}
}
