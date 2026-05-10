package restapi

import (
	"net/http"

	"github.com/71g3pf4c3/gophermart/internal/controller/restapi/middleware"
	"github.com/71g3pf4c3/gophermart/pkg/jwt"
	"github.com/71g3pf4c3/gophermart/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// NewRouter configures HTTP routes for the Gophermart service.
func NewRouter(app *fiber.App, jwtManager *jwt.Manager, l logger.Interface, userService User, orderService Order, balanceService Balance) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Get("/ping", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	api := app.Group("/api")
	user := api.Group("/user")

	uh := newUserHandler(userService)
	user.Post("/register", uh.Register)
	user.Post("/login", uh.Login)

	auth := user.Group("", middleware.Auth(jwtManager))

	oh := newOrderHandler(orderService)
	auth.Post("/orders", oh.Upload)
	auth.Get("/orders", oh.List)

	bh := newBalanceHandler(balanceService)
	auth.Get("/balance", bh.GetBalance)
	auth.Post("/balance/withdraw", bh.Withdraw)
	auth.Get("/withdrawals", bh.GetWithdrawals)
}
