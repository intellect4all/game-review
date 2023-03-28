package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
)

func router(ctx context.Context, app fiber.Router, handler *Handler, middleware auth.Middleware) error {

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion + "/reviews")

	app.Use(middleware.AuthMiddleware(allowAllAuthenticated))

	app.Post("/add", HandleAddReview(handler, ctx))

	app.Put("/:id", HandleUpdateReview(handler, ctx))

	app.Get("/:id", HandleGetReview(handler, ctx))

	app.Get("/game/:gameId", HandleGetReviewsForFame(handler, ctx))

	app.Get("/user/:userId", HandleGetReviewsForUser(handler, ctx))

	app.Delete("/:id", HandleDeleteReview(handler, ctx))

	app.Post("/upVote/:id", HandleVoteReview(handler, ctx, true))

	app.Post("/downVote/:id", HandleVoteReview(handler, ctx, false))

	app.Use(middleware.AuthMiddleware(isAdminOrModerator))

	app.Get("/flagged", HandleGetFlaggedReviews(handler, ctx))

	app.Post("/flag/:id", HandleFlagReview(handler, ctx))

	return nil
}

func isAdminOrModerator(claims *auth.JwtClaims) (string, bool) {

	if claims.Role == "admin" || claims.Role == "moderator" {
		return "", true
	}

	return "", false
}

func allowAllAuthenticated(claims *auth.JwtClaims) (string, bool) {
	return "", true
}
