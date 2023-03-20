package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func authRouter(ctx context.Context, app fiber.Router, handler *AuthHandler) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion)

	app.Post("/login", HandleLogin(handler, ctx))

	app.Post("/signup", HandleSignUp(handler, ctx))

	app.Post("/verify-account/init/:email", HandleVerifyAccountInit(handler, ctx))

	app.Post("/verify-account/verify", HandleVerifyAccount(handler, ctx))

	app.Post("/verify-account/resend/:email", HandleVerifyAccountResend(handler, ctx))

	app.Post("/forgot-password/init/:email", HandleForgetPassword(handler, ctx))

	app.Post("/forgot-password/reset", HandlePasswordReset(handler, ctx))

	app.Post("/forgot-password/resend/:email", HandleForgetPasswordResend(handler, ctx))

	return nil
}
