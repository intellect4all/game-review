package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
)

func router(ctx context.Context, app fiber.Router, handler *GameHandler, middleware auth.AuthMiddleware) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion + "/games")

	app.Use(middleware.GetMiddleWare(check))

	app.Post("/genres/add", HandleAddGenre(handler, ctx))

	return nil
}

func check(claims *auth.JwtClaims) (string, bool) {

	if claims.Role == "admin" || claims.Role == "moderator" {
		return "", true
	}

	return "You are not authorized to perform this action", false
}
