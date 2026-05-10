package restapi

import (
	"context"
	"errors"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// Balance is the usecase interface for balance operations.
type Balance interface {
	GetBalance(ctx context.Context, userID string) (entity.Balance, error)
	Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}

type balanceHandler struct {
	uc Balance
}

func newBalanceHandler(uc Balance) *balanceHandler {
	return &balanceHandler{uc: uc}
}

// GetBalance handles GET /api/user/balance.
func (h *balanceHandler) GetBalance(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	b, err := h.uc.GetBalance(c.UserContext(), userID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(b)
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// Withdraw handles POST /api/user/balance/withdraw.
func (h *balanceHandler) Withdraw(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req withdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Order == "" || req.Sum <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := h.uc.Withdraw(c.UserContext(), userID, req.Order, req.Sum)
	switch {
	case err == nil:
		return c.SendStatus(fiber.StatusOK)
	case errors.Is(err, entity.ErrInvalidOrderNumber):
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	case errors.Is(err, entity.ErrInsufficientFunds):
		return c.SendStatus(fiber.StatusPaymentRequired)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

// GetWithdrawals handles GET /api/user/withdrawals.
func (h *balanceHandler) GetWithdrawals(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	withdrawals, err := h.uc.GetWithdrawals(c.UserContext(), userID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(withdrawals) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.JSON(withdrawals)
}
