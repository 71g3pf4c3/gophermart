package restapi

import (
	"context"
	"errors"
	"strings"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type User interface {
	Register(ctx context.Context, login, password string) (entity.User, string, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type UserHandler struct {
	u User
}

type userAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func newUserHandler(u User) *UserHandler {
	return &UserHandler{u: u}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req userAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, token, err := h.u.Register(c.UserContext(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			return c.SendStatus(fiber.StatusConflict)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	c.Set(fiber.HeaderAuthorization, "Bearer "+token)

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req userAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := h.u.Login(c.UserContext(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidCredentials):
			return c.SendStatus(fiber.StatusUnauthorized)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	c.Set(fiber.HeaderAuthorization, "Bearer "+token)

	return c.SendStatus(fiber.StatusOK)
}
