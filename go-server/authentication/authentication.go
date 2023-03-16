package authentication

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

func Register(mongoClient *mongo.Client, ctx context.Context, app fiber.Router) error {
	authRepo := NewMongoRepository(mongoClient)
	jwtHelper := NewJWTHelper()
	authService := NewService(authRepo, jwtHelper)
	authHandler := NewHandler(authService)

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion)

	app.Post("/login", func(c *fiber.Ctx) error {
		return authHandler.Login(ctx, c)
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		return authHandler.Signup(ctx, c)
	})

	app.Post("verify-account/init", func(c *fiber.Ctx) error {
		return authHandler.InitAccountVerification(ctx, c)
	})

	app.Post("verify-account/verify", func(c *fiber.Ctx) error {
		return authHandler.VerifyAccount(ctx, c)
	})

	return nil
}

type AuthHandler struct {
	authService *Service
}

func NewHandler(authService *Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(ctx context.Context, c *fiber.Ctx) error {
	var userCredential UserCredential
	err := c.BodyParser(&userCredential)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	jwt, err := a.authService.AuthenticateUser(ctx, &userCredential)
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

	var userCredential *UserCredential
	err := c.BodyParser(&userCredential)

	fmt.Printf("password: %s, id: %s, role, %s", userCredential.Password, userCredential.Id, userCredential.Role)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	setDefaultValues(userCredential)

	err = a.authService.CreateUser(ctx, userCredential)

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

func (a *AuthHandler) InitAccountVerification(ctx context.Context, c *fiber.Ctx) error {

	type Email struct {
		Email string `json:"email" bson:"email"`
	}

	var email Email

	err := c.BodyParser(&email)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	userId := UserID(email.Email)
	fmt.Printf("email: %s, userID: %s", email, userId)
	otp, err := a.authService.CreateVerificationOTP(ctx, &userId)
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

	return c.Status(fiber.StatusCreated).JSON(OTPCreationSuccessResponse(otp, email.Email))
}

func (a *AuthHandler) VerifyAccount(ctx context.Context, c *fiber.Ctx) error {

	var verifyRequest VerifyUserRequest

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

func setDefaultValues(credential *UserCredential) {
	credential.IsActive = strings.ToLower(credential.Role) == "user"
	credential.CreatedAt = time.Now()
	credential.IsVerified = false
}

func getFiberStatusFromError(err error) int {
	switch err {
	case ErrUserAlreadyExists:
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}

}

//func BookSuccessResponse(data *entities.Book) *fiber.Map {
//	book := Book{
//		ID:     data.ID,
//		Title:  data.Title,
//		Author: data.Author,
//	}
//	return &fiber.Map{
//		"status": true,
//		"data":   book,
//		"error":  nil,
//	}
//}
//
//// BooksSuccessResponse is the list SuccessResponse that will be passed in the response by Handler
//func BooksSuccessResponse(data *[]Book) *fiber.Map {
//	return &fiber.Map{
//		"status": true,
//		"data":   data,
//		"error":  nil,
//	}
//}
//
//// BookErrorResponse is the ErrorResponse that will be passed in the response by Handler
//func BookErrorResponse(err error) *fiber.Map {
//	return &fiber.Map{
//		"status": false,
//		"data":   "",
//		"error":  err.Error(),
//	}
//}
