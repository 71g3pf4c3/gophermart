package restapi

import (
	"context"
	"errors"
	"strings"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// Order is the usecase interface for order operations.
type Order interface {
	Upload(ctx context.Context, userID, number string) error
	List(ctx context.Context, userID string) ([]entity.Order, error)
}

type orderHandler struct {
	uc Order
}

func newOrderHandler(uc Order) *orderHandler {
	return &orderHandler{uc: uc}
}

// Upload handles POST /api/user/orders.
func (h *orderHandler) Upload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	number := strings.TrimSpace(string(c.Body()))
	if number == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := h.uc.Upload(c.UserContext(), userID, number)
	switch {
	case err == nil:
		return c.SendStatus(fiber.StatusAccepted)
	case errors.Is(err, entity.ErrOrderOwnedBySelf):
		return c.SendStatus(fiber.StatusOK)
	case errors.Is(err, entity.ErrOrderOwnedByOther):
		return c.SendStatus(fiber.StatusConflict)
	case errors.Is(err, entity.ErrInvalidOrderNumber):
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

// List handles GET /api/user/orders.
func (h *orderHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orders, err := h.uc.List(c.UserContext(), userID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(orders) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.JSON(orders)
}
