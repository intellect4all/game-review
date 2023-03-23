package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
)

func router(ctx context.Context, app fiber.Router, handler *GameHandler, middleware auth.Middleware) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion + "/games")

	app.Use(middleware.AuthMiddleware(allowAllAuthenticated))

	app.Get("/genres", HandleGetGenres(handler, ctx))

	app.Get("/genres/:slug", HandleGetGenre(handler, ctx))

	app.Use(middleware.AuthMiddleware(check))

	app.Post("/genres/add", HandleAddGenre(handler, ctx))

	app.Post("/genres/update", HandleUpdateGenre(handler, ctx))

	app.Use(middleware.AuthMiddleware(adminOnlyPermission))

	app.Delete("/genres/:slug", HandleDeleteGenre(handler, ctx))

	return nil
}

func adminOnlyPermission(claims *auth.JwtClaims) (string, bool) {

	if claims.Role != "admin" {
		return "", false
	}

	return "", true
}

func allowAllAuthenticated(claims *auth.JwtClaims) (string, bool) {
	return "", true
}

func check(claims *auth.JwtClaims) (string, bool) {

	if claims.Role == "admin" || claims.Role == "moderator" {
		return "", true
	}

	return "You are not authorized to perform this action", false
}
