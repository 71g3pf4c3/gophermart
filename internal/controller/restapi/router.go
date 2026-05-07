package restapi

import (
	"net/http"

	"github.com/71g3pf4c3/gophermart/internal/controller/restapi/middleware"
	"github.com/71g3pf4c3/gophermart/pkg/jwt"
	"github.com/71g3pf4c3/gophermart/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// NewRouter configures HTTP routes for the Gophermart service.
func NewRouter(app *fiber.App, jwtManager *jwt.Manager, l logger.Interface, userService User) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	app.Get("/ping", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	api := app.Group("/api")
	user := api.Group("/user")
	us := newUserHandler(userService)
	user.Post("/register", us.Register)
	user.Post("/login", us.Login)

	authenticated := user.Group("", middleware.Auth(jwtManager))
	authenticated.Get("/users", notImplemented)
	authenticated.Post("/orders", notImplemented)
	authenticated.Get("/orders", notImplemented)
	authenticated.Get("/balance", notImplemented)
	authenticated.Post("/balance/withdraw", notImplemented)
	authenticated.Get("/withdrawals", notImplemented)
}

func notImplemented(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
}
