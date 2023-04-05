package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
)

// HandleAddReview  godoc
//
// @Security BearerAuth
//
// @Summary Add a review
// @Description Add a review
// @Tags Reviews
// @ID addReview
// @Accept json
// @Produce json
//
// @Param addReview body reviews.AddReviewRequest true "addReview request"
//
// @Success 201 {object} main.JSONResult{data=reviews.AddReviewRes} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Game not found"
// @Failure 409 {object} main.JSONErrorRes "Review already exists"
// @Router /api/v1/reviews/add [post]
func HandleAddReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.AddReview(ctx, c)
	}
}

// HandleUpdateReview  godoc
//
// @Security BearerAuth
//
// @Summary Update a review
// @Description Update a review
// @Tags Reviews
// @ID updateReview
// @Accept json
// @Produce json
//
// @Param addReview body reviews.AddReviewRequest true "update review request"
// @Param reviewId path string true "review id"
//
// @Success 201 {object} main.JSONResult{data=reviews.AddReviewRes} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId} [put]
func HandleUpdateReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.UpdateReview(ctx, c)
	}
}

// HandleGetReview godoc
//
// @Security BearerAuth
//
// @Summary Get a review
// @Description Get a review
// @Tags Reviews
// @ID getReview
// @Accept json
// @Produce json
//
// @Param reviewId path string true "review id"
//
// @Success 200 {object} main.JSONResult{data=reviews.ReviewResponse} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId} [get]
func HandleGetReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReview(ctx, c)
	}
}

// HandleGetReviewsForGame godoc
//
// @Security BearerAuth
//
// @Summary Get reviews for a game
// @Description Get reviews for a game
// @Tags Reviews
// @ID getReviewsForGame
// @Accept json
// @Produce json
//
// @Param getGameReviews body reviews.GetReviewsForGame	 true "getGameReviews request"
// @Param gameId path string true "game id"
//
// @Success 200 {object} main.JSONResult{data=reviews.PaginatedResponse[ReviewResponse]} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Game not found"
// @Router /api/v1/reviews/game/{gameId} [get]
func HandleGetReviewsForGame(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReviewsForGame(ctx, c)
	}
}

// HandleGetReviewsForUser  godoc
//
// @Security BearerAuth
//
// @Summary Get all reviews for a user
// @Description Get all reviews for a user
// @Tags Reviews
// @ID getReviewsForUser
// @Accept json
// @Produce json
//
// @Param getGameReviews body reviews.GetReviewsForGame	 true "getGameReviews request"
// @Param userId path string true "user id"
//
// @Success 200 {object} main.JSONResult{data=reviews.PaginatedResponse[ReviewResponse]} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "User not found"
// @Router /api/v1/reviews/user/{userId} [get]
func HandleGetReviewsForUser(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReviewsForUser(ctx, c)
	}
}

// HandleDeleteReview  godoc
//
// @Security BearerAuth
//
// @Summary Delete a review
// @Description Delete a review
// @Tags Reviews
// @ID deleteReview
// @Accept json
// @Produce json
//
// @Param reviewId path string true "review id"
//
// @Success 200 {object} main.JSONResult{data=reviews.AddReviewRes} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId} [delete]
func HandleDeleteReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.DeleteReview(ctx, c)
	}
}

func HandleVoteReview(handler *Handler, ctx context.Context, shouldUpvote bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		if shouldUpvote {
			return upvote(ctx, c, handler)
		}
		return downVote(ctx, c, handler)
	}
}

// HandleUpvoteReview godoc
//
// @Security BearerAuth
//
// @Summary Upvote a review
// @Description Upvote a review
// @Tags Reviews
// @ID upvoteReview
// @Accept json
// @Produce json
//
// @Param reviewId path string true "review id"
//
// @Success 200 {object} main.JSONResult{data=string} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId}/upvote [post]
func upvote(ctx context.Context, c *fiber.Ctx, handler *Handler) error {
	return handler.VoteReview(ctx, c, true)
}

// HandleDownvote godoc
//
// @Security BearerAuth
//
// @Summary Upvote a review
// @Description Upvote a review
// @Tags Reviews
// @ID downvoteReview
// @Accept json
// @Produce json
//
// @Param reviewId path string true "review id"
//
// @Success 200 {object} main.JSONResult{data=string} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId}/downvote [post]
func downVote(ctx context.Context, c *fiber.Ctx, handler *Handler) error {
	return handler.VoteReview(ctx, c, false)
}

// HandleGetFlaggedReviews godoc
//
// @Security BearerAuth
//
// @Summary Get flagged reviews
// @Description Get flagged reviews
// @Tags Reviews
// @ID getFlaggedReviews
// @Accept json
// @Produce json
//
// @Param getGameReviews body reviews.GetFlaggedReviewsRequest	 true "getFlaggedRequest request"
//
// @Success 200 {object} main.JSONResult{data=reviews.PaginatedResponse[Review]} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Router /api/v1/reviews/flagged [get]
func HandleGetFlaggedReviews(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetFlaggedReviews(ctx, c)
	}
}

// HandleFlagReview godoc
//
// @Security BearerAuth
//
// @Summary Flag a review
// @Description Flag a review
// @Tags Reviews
// @ID flagReview
// @Accept json
// @Produce json
//
// @Param reviewId path string true "review id"
//
// @Success 200 {object} main.JSONResult{data=string} "Success"
// @Failure 400 {object} main.JSONErrorRes "Bad request"
// @Failure 404 {object} main.JSONErrorRes "Review not found"
// @Router /api/v1/reviews/{reviewId}/flag [post]
func HandleFlagReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.FlagReview(ctx, c)
	}
}

func GetNewContext(ctx context.Context, c *fiber.Ctx) context.Context {
	userId := c.Locals("userId").(string)
	role := c.Locals("role").(string)

	log.Println("userId", userId)
	log.Println("role", role)
	log.Println("userId", userId)
	log.Println("role", role)
	ctx = context.WithValue(ctx, "userId", userId)
	ctx = context.WithValue(ctx, "role", role)
	return ctx
}
