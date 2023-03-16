package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func authRouter(ctx context.Context, app fiber.Router, handler *AuthHandler) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion)

	app.Post("/login", func(c *fiber.Ctx) error {
		return handler.Login(ctx, c)
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		return handler.Signup(ctx, c)
	})

	app.Post("/verify-account/init", func(c *fiber.Ctx) error {
		return handler.InitAccountVerification(ctx, c)
	})

	app.Post("/verify-account/verify", func(c *fiber.Ctx) error {
		return handler.VerifyAccount(ctx, c)
	})

	app.Post("/verify-account/resend", func(c *fiber.Ctx) error {
		return handler.InitAccountVerification(ctx, c)
	})

	app.Post("/forgot-password/init", func(c *fiber.Ctx) error {
		return handler.InitForgotPassword(ctx, c)
	})

	app.Post("/forgot-password/reset", func(c *fiber.Ctx) error {
		return handler.VerifyForgotPassword(ctx, c)
	})

	app.Post("/forgot-password/resend", func(c *fiber.Ctx) error {
		return handler.InitForgotPassword(ctx, c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return handler.Login(ctx, c)
	})

	return nil
}
