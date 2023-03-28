package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func HandleAddReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.AddReview(ctx, c)
	}
}

func HandleUpdateReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.UpdateReview(ctx, c)
	}
}

func HandleGetReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReview(ctx, c)
	}
}

func HandleGetReviewsForFame(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReviewsForGame(ctx, c)
	}
}

func HandleGetReviewsForUser(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetReviewsForUser(ctx, c)
	}
}

func HandleDeleteReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.DeleteReview(ctx, c)
	}
}

func HandleVoteReview(handler *Handler, ctx context.Context, shouldUpvote bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.VoteReview(ctx, c, shouldUpvote)
	}
}

func HandleGetFlaggedReviews(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetFlaggedReviews(ctx, c)
	}
}

func HandleFlagReview(handler *Handler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.FlagReview(ctx, c)
	}
}

func GetNewContext(ctx context.Context, c *fiber.Ctx) context.Context {
	userId := c.Locals("userId").(string)
	role := c.Locals("role").(string)
	ctx = context.WithValue(ctx, "userId", userId)
	ctx = context.WithValue(ctx, "role", role)
	return ctx
}
