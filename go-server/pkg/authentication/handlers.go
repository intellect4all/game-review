package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

type AuthHandler struct {
	authService *AuthService
}

func NewHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(ctx context.Context, c *fiber.Ctx) error {
	var req LoginRequest
	err := c.BodyParser(&req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	userData, jwt, err := a.authService.AuthenticateUser(ctx, &req)
	if err != nil {
		var status int
		var message string
		if err == ErrAccountInactive {
			status = fiber.StatusUpgradeRequired
			message = "Account is marked inactive. Contact Support"
		} else if err == ErrInvalidCredentials {
			status = fiber.StatusBadRequest
			message = "Invalid credentials"
		} else if err == ErrUserNotFound {
			status = fiber.StatusNotFound
			message = "User not found"
		} else {
			status = fiber.StatusInternalServerError
			message = "Something went wrong"
		}

		return c.Status(status).JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetLoginSuccessResponse(jwt, userData))

}

func (a *AuthHandler) Signup(ctx context.Context, c *fiber.Ctx) error {

	var req *SignUpRequest
	err := c.BodyParser(&req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	err = a.authService.CreateUser(ctx, *req)

	if err != nil {
		status := 500
		if err == ErrUserAlreadyExists {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "Error creating user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func (a *AuthHandler) InitAccountVerification(ctx context.Context, c *fiber.Ctx) error {

	userId, err := extractEmailFromPathParams(c)

	if err != nil {
		return err
	}

	otpID, err := a.authService.CreateVerificationOTP(ctx, userId)
	if err != nil {
		status := 0
		switch err {
		case ErrUserAlreadyVerified:
			status = fiber.StatusConflict
		case ErrUserNotFound:
			status = fiber.StatusNotFound
		default:
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "Error creating OTP",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(GetOTPCreationResponse(&otpID, &userId))
}

func (a *AuthHandler) VerifyAccount(ctx context.Context, c *fiber.Ctx) error {

	var verifyRequest VerifyAccountRequest

	err := c.BodyParser(&verifyRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"error":   ErrInvalidRequest.Error(),
		})
	}

	err = a.authService.VerifyUser(ctx, verifyRequest)

	if err != nil {
		status := 0
		switch err {
		case ErrUserAlreadyVerified:
			status = fiber.StatusConflict
		case ErrUserNotFound:
			status = fiber.StatusNotFound
		case ErrInvalidOTP, ErrOTPExpired, ErrOTPUsed:
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "Error verifying account",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account verified successfully",
	})

}

func (a *AuthHandler) InitForgotPassword(ctx context.Context, c *fiber.Ctx) error {
	userID, err := extractEmailFromPathParams(c)

	if err != nil {
		return err
	}

	otpID, err := a.authService.InitForgotPassword(ctx, userID)

	if err != nil {
		status := 0
		if err == ErrUserNotFound {
			status = fiber.StatusNotFound
		} else {
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "Error initiating forgot password",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(GetOTPCreationResponse(&otpID, &userID))
}

func (a *AuthHandler) ResetPassword(ctx context.Context, c *fiber.Ctx) error {
	var verifyRequest ForgetAndResetPasswordRequest

	err := c.BodyParser(&verifyRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"error":   ErrInvalidRequest.Error(),
		})
	}

	err = a.authService.ChangePassword(ctx, verifyRequest)

	if err != nil {
		status := 0
		switch err {
		case ErrUserNotFound:
			status = fiber.StatusNotFound
		case ErrInvalidOTP, ErrOTPExpired, ErrOTPUsed, ErrPasswordMismatch:
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "Error resetting password",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
	})

}

func extractEmailFromPathParams(c *fiber.Ctx) (string, error) {
	email := struct {
		Email string `params:"email"`
	}{}

	err := c.ParamsParser(&email)

	emailStr, err := url.QueryUnescape(email.Email)

	if err != nil {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   ErrInvalidRequest,
		})
	}

	return emailStr, nil
}
