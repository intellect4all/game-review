package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func authRouter(ctx context.Context, app fiber.Router, handler *AuthHandler) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion)

	app.Post("/account/login", HandleLogin(handler, ctx))

	app.Post("/account/signup", HandleSignUp(handler, ctx))

	app.Post("/account/init-verification/:email", HandleVerifyAccountInit(handler, ctx))

	app.Post("/account/verify-email", HandleVerifyAccount(handler, ctx))

	app.Post("/account/verify-email/resend/:email", HandleVerifyAccountResend(handler, ctx))

	app.Post("/account/forgot-password/init/:email", HandleForgetPassword(handler, ctx))

	app.Post("/account/forgot-password/reset", HandlePasswordReset(handler, ctx))

	app.Post("/account/forgot-password/resend/:email", HandleForgetPasswordResend(handler, ctx))

	return nil
}
