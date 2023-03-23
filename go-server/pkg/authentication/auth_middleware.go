package authentication

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

type AuthMiddleware interface {
	GetMiddleWare(authCheck func(claims *JwtClaims) (string, bool)) interface{}
}

type AuthMiddlewareImpl struct {
	jwtHelper JWTHelper
}

func NewAuthMiddleware(jwtHelper JWTHelper) *AuthMiddlewareImpl {
	return &AuthMiddlewareImpl{
		jwtHelper: jwtHelper,
	}
}

func (a *AuthMiddlewareImpl) GetMiddleWare(authCheck func(claims *JwtClaims) (string, bool)) interface{} {
	return func(c *fiber.Ctx) error {
		log.Println("AuthMiddleware called")
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("authToken is empty")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header not found",
			})
		}

		authToken := strings.Split(authHeader, "Bearer ")[1]
		if authToken == "" {
			log.Println("authToken is empty")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization token not found",
			})
		}

		log.Println("authToken: ", authToken)
		var claims, err = a.jwtHelper.ValidateJWT(AuthenticatedUserJWT(authToken))
		log.Println("Claims validated")

		log.Println("claims: ", claims)
		if err != nil {
			log.Println("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		log.Println("claims: ", claims)

		if message, ok := authCheck(claims); !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": message,
				"error":   "Unauthorized",
			})
		}

		return c.Next()
	}
}
