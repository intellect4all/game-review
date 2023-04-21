package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
	"log"
)

func router(ctx context.Context, app fiber.Router, handler *Handler, middleware auth.Middleware) error {

	log.Println("Registering reviews routes")

	apiVersion := ctx.Value("apiVersion").(string)
	app = app.Group(apiVersion + "/reviews")

	app.Use(middleware.AuthMiddleware(allowAllAuthenticated))

	app.Post("/add", HandleAddReview(handler, ctx))

	app.Put("/:id", HandleUpdateReview(handler, ctx))

	app.Get("/flagged/", HandleGetFlaggedReviews(handler, ctx))

	app.Get("/locations/", HandleGetReviewsLocations(handler, ctx))

	app.Get("/:id/", HandleGetReview(handler, ctx))

	app.Get("/game/:gameId/", HandleGetReviewsForGame(handler, ctx))

	app.Get("/user/:userId/", HandleGetReviewsForUser(handler, ctx))

	app.Delete("/:id/", HandleDeleteReview(handler, ctx))

	app.Post("/:id/upvote/", HandleVoteReview(handler, ctx, true))

	app.Post("/:id/downvote/", HandleVoteReview(handler, ctx, false))

	app.Use(middleware.AuthMiddleware(isAdminOrModerator))
	app.Post("/:id/unflag/", HandleUnflagReview(handler, ctx))

	app.Post("/:id/flag/", HandleFlagReview(handler, ctx))

	return nil
}

func isAdminOrModerator(claims *auth.JwtClaims) (string, bool) {
	log.Println("isAdminOrModerator")

	if claims.Role == "admin" || claims.Role == "moderator" {
		return "", true
	}

	return "", false
}

func allowAllAuthenticated(claims *auth.JwtClaims) (string, bool) {
	log.Println("allowAllAuthenticated")
	return "", true
}
