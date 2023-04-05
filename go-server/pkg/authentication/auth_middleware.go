package authentication

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

type Middleware interface {
	AuthMiddleware(authCheck func(claims *JwtClaims) (string, bool)) interface{}
}

type AuthMiddlewareImpl struct {
	jwtHelper JWTHelper
}

func NewAuthMiddleware(jwtHelper JWTHelper) *AuthMiddlewareImpl {
	return &AuthMiddlewareImpl{
		jwtHelper: jwtHelper,
	}
}

func (a *AuthMiddlewareImpl) AuthMiddleware(authCheck func(claims *JwtClaims) (string, bool)) interface{} {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header not found",
			})
		}

		authToken := strings.Split(authHeader, "Bearer ")[1]
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization token not found",
			})
		}

		var claims, err = a.jwtHelper.ValidateJWT(AuthenticatedUserJWT(authToken))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		if message, ok := authCheck(claims); !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": message,
				"error":   "Unauthorized",
			})
		}

		c.Locals("userId", claims.Id)
		c.Locals("role", claims.Role)

		log.Println("user validated")

		return c.Next()
	}
}

func (a *AuthMiddlewareImpl) RouteGuard(authCheck func(claims *JwtClaims) (string, bool)) interface{} {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header not found",
			})
		}

		authToken := strings.Split(authHeader, "Bearer ")[1]
		if authToken == "" {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization token not found",
			})
		}

		var claims, err = a.jwtHelper.ValidateJWT(AuthenticatedUserJWT(authToken))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		if message, ok := authCheck(claims); !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": message,
				"error":   "Unauthorized",
			})
		}

		return c.Next()
	}
}
