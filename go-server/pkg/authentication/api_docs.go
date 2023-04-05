package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

// HandleLogin godoc
//
//	@Summary		Login endpoint for all users
//	@Description	Returns a signed JSON Web Token that can be used to talk to secured endpoints
//	@Tags			Account
//	@ID				login
//	@Accept			json
//	@Produce		json
//
//	@Param			loginRequest	body		authentication.LoginRequest 	true			"login request"
//
//	@Success		200				{object}	main.JSONResult{data=authentication.LoginDTO}	"success"
//	@Failure		426				{object}	main.JSONErrorRes											"Account is inactive"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		404				{object}	main.JSONErrorRes											"User not found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/login [post]
func HandleLogin(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.Login(ctx, c)
	}
}

// HandleSignUp  godoc
//
//	@Summary		Signup endpoint for all users and moderators
//	@Description	Creates a new User/Moderator on the system. The Moderator will need to be manually activated by an existing admin
//	@Tags			Account
//	@ID				signup
//	@Accept			json
//	@Produce		json
//
//	@Param			signUpRequest	body		authentication.SignUpRequest 	true			"signup request"
//
//	@Success		200				{object}	main.JSONResult{data=string}	"Success"
//
// @Success		207				{object}	main.JSONResult{data=string}	"Success"
//
//	@Failure		409				{object}	main.JSONErrorRes					"User already exists"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/signup [post]
func HandleSignUp(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.Signup(ctx, c)
	}
}

// HandleVerifyAccountInit  godoc
//
//	@Summary		Initiate user email verification
//	@Description	An otp code is sent to the email if the user account existed.
//
// @Description 	An otp ID is returned, which must submitted alongside the otpCode sent to the mail to the VerifyAccount endpoint.
//
//	@Tags			Account
//	@ID				initVerifyAccount
//	@Accept			json
//	@Produce		json
//
//	@Param			email	path		string 	true			"Email address"
//
//	@Success		200				{object}	main.JSONResult{data=authentication.OTPCreationSuccessResponse}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		409				{object}	main.JSONErrorRes											"User already verified"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/init-verification/{email} [post]
func HandleVerifyAccountInit(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.InitAccountVerification(ctx, c)
	}
}

// HandleVerifyAccount   godoc
//
//	@Summary		Complete account verification
//	@Description	The Endpoint verifies the user account using the otp code sent to the user's email
//
//	@Tags			Account
//	@ID				VerifyAccount
//	@Accept			json
//	@Produce		json
//
//	@Param			verifyAccountRequest	body		authentication.VerifyAccountRequest 	true "OtpID data"
//
//	@Success		200				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		409				{object}	main.JSONErrorRes											"User already verified"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/verify-email [post]
func HandleVerifyAccount(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.VerifyAccount(ctx, c)
	}
}

// HandleVerifyAccountResend   godoc
//
//	@Summary		Resend OTP code for account verification
//	@Description	An otp code is sent to the email if the user account existed.
//
// @Description 	An otp ID is returned, which must submitted alongside the otpCode sent to the mail to the VerifyAccount endpoint.
//
//	@Tags			Account
//	@ID				ResendVerificationCode
//	@Accept			json
//	@Produce		json
//
//	@Param			email	path		string 	true			"Email address"
//
//	@Success		200				{object}	main.JSONResult{data=authentication.OTPCreationSuccessResponse}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		409				{object}	main.JSONErrorRes											"User already verified"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/verify-email/resend/{email} [post]
func HandleVerifyAccountResend(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.InitAccountVerification(ctx, c)
	}
}

// HandleForgetPassword  godoc
//
//	@Summary		Forget password endpoint
//	@Description	An otp code is sent to the email if the user account existed.
//
// @Description 	An otp ID is returned, which must submitted alongside the otpCode sent to the mail to the Forgot password reset endpoint.
//
//	@Tags			Account
//	@ID				ForgotPasswordInit
//	@Accept			json
//	@Produce		json
//
//	@Param			email	path		string 	true			"Email address"
//
//	@Success		200				{object}	main.JSONResult{data=authentication.OTPCreationSuccessResponse}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/forgot-password/init/{email} [post]
func HandleForgetPassword(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.InitForgotPassword(ctx, c)
	}
}

// HandlePasswordReset   godoc
//
//	@Summary		Complete Forget password reset
//	@Description	The Endpoint resets the user password using the otp code sent to the user's email
//
//	@Tags			Account
//	@ID				ResetPassword
//	@Accept			json
//	@Produce		json
//
//	@Param			ForgetAndResetPasswordRequest	body		authentication.ForgetAndResetPasswordRequest 	true "OtpID data"
//
//	@Success		200				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		409				{object}	main.JSONErrorRes											"User already verified"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/forgot-password/reset [post]
func HandlePasswordReset(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.ResetPassword(ctx, c)
	}
}

// HandleForgetPasswordResend   godoc
//
//	@Summary		Resend OTP code for Forget password
//	@Description	An otp code is sent to the email if the user account existed.
//
// @Description 	An otp ID is returned, which must submitted alongside the otpCode sent to the mail to the Forgot password reset endpoint.
//
//	@Tags			Account
//	@ID				ResendForgetPasswordCode
//	@Accept			json
//	@Produce		json
//
//	@Param			email	path		string 	true			"Email address"
//
//	@Success		200				{object}	main.JSONResult{data=authentication.OTPCreationSuccessResponse}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"No Account Found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/account/forgot-password/resend/{email} [post]
func HandleForgetPasswordResend(handler *AuthHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.InitForgotPassword(ctx, c)
	}
}
