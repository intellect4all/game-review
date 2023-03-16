package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *Service
}

func NewHandler(authService *Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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

	jwt, err := a.authService.AuthenticateUser(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    map[string]interface{}{"jwt": jwt},
	})

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

	err = a.authService.CreateUser(ctx, req)

	if err != nil {
		status := getFiberStatusFromError(err)
		return c.Status(status).JSON(fiber.Map{
			"message": "Error creating user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

type Email struct {
	Email string `json:"email" bson:"email"`
}

func (a *AuthHandler) InitAccountVerification(ctx context.Context, c *fiber.Ctx) error {

	userId, err := extractUserIDFromRequest(c)

	if err != nil {
		return err
	}

	otp, err := a.authService.CreateVerificationOTP(ctx, userId)
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

	return c.Status(fiber.StatusCreated).JSON(OTPCreationSuccessResponse(&otp, (*string)(userId)))
}

func (a *AuthHandler) VerifyAccount(ctx context.Context, c *fiber.Ctx) error {

	var verifyRequest OtpCheckData

	err := c.BodyParser(&verifyRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"error":   ErrInvalidRequest.Error(),
		})
	}

	err = a.authService.VerifyUser(ctx, &verifyRequest)

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
	userID, err := extractUserIDFromRequest(c)

	if err != nil {
		return err
	}

	otp, err := a.authService.InitForgotPassword(ctx, userID)

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

	return c.Status(fiber.StatusCreated).JSON(OTPCreationSuccessResponse(&otp, (*string)(userID)))
}

func (a *AuthHandler) VerifyForgotPassword(ctx context.Context, c *fiber.Ctx) error {
	var verifyRequest ForgetAndResetPasswordRequest

	err := c.BodyParser(&verifyRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"error":   ErrInvalidRequest.Error(),
		})
	}

	err = a.authService.ChangePassword(ctx, &verifyRequest)

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

func extractUserIDFromRequest(c *fiber.Ctx) (*UserID, error) {
	var email Email

	err := c.BodyParser(&email)

	var userID UserID

	if err != nil {
		return &userID, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	userID = UserID(email.Email)

	return &userID, nil
}

func getFiberStatusFromError(err error) int {
	switch err {
	case ErrUserAlreadyExists:
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}

}
