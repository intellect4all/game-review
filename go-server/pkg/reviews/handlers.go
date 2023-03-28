package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

type AddReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,number,gte=0,lte2"`
	Comment string `json:"comment" validate:"required,min=5,max=2000"`
	GameId  string `json:"gameId" validate:"required"`
}

func (h *Handler) AddReview(ctx context.Context, c *fiber.Ctx) error {
	var req AddReviewRequest

	err := c.BodyParser(&req)

	if err != nil {
		return AddReviewErrorResponse(c, ErrBadRequest)
	}
	userId := c.Locals("userId").(string)
	review := &AddReview{
		Rating:  req.Rating,
		Comment: req.Comment,
		GameId:  req.GameId,
		UserId:  userId,
	}

	id, err := h.Service.addReview(ctx, review)

	if err != nil {
		return AddReviewErrorResponse(c, err)
	}

	return AddReviewSuccessResp(c, id)
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,number,gte=0,lte2"`
	Comment string `json:"comment" validate:"required,min=5,max=2000"`
	GameId  string `json:"gameId" validate:"required"`
}

func (h *Handler) UpdateReview(ctx context.Context, c *fiber.Ctx) error {
	var req AddReviewRequest

	err := c.BodyParser(&req)

	if err != nil {
		return UpdateReviewErrorResponse(c, ErrBadRequest)
	}

	id := c.Params("id")

	review := &AddReview{
		Rating:  req.Rating,
		Comment: req.Comment,
		GameId:  req.GameId,
	}

	err = h.Service.updateReview(ctx, id, review)

	if err != nil {
		return UpdateReviewErrorResponse(c, err)
	}

	return UpdateReviewSuccessResp(c)
}

func (h *Handler) GetReview(ctx context.Context, c *fiber.Ctx) error {
	id := c.Params("id")

	review, err := h.Service.getReview(ctx, id)

	if err != nil {
		return GetReviewErrorResponse(c, err)
	}

	return GetReviewSuccessResp(c, *review)
}

func (h *Handler) GetReviewsForGame(ctx context.Context, c *fiber.Ctx) error {
	return getReview(ctx, c, true, h.Service.getReviewsForGame)
}

func (h *Handler) GetReviewsForUser(ctx context.Context, c *fiber.Ctx) error {
	return getReview(ctx, c, false, h.Service.getReviewsForUser)
}

func (h *Handler) DeleteReview(ctx context.Context, c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.Service.deleteReview(ctx, id)

	if err != nil {
		return GetReviewErrorResponse(c, err)
	}

	return DeleteReviewSuccessResp(c)

}

func (h *Handler) VoteReview(ctx context.Context, c *fiber.Ctx, upvote bool) error {

	id := c.Params("id")

	err := h.Service.voteReview(ctx, id, upvote)

	if err != nil {
		return VoteReviewErrorResponse(c, err)
	}

	return VoteReviewSuccessResp(c)

}

type GetFlaggedReviewsRequest struct {
	GameId string `query:"gameId"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

func (h *Handler) GetFlaggedReviews(ctx context.Context, c *fiber.Ctx) error {
	var req GetFlaggedReviewsRequest

	err := c.QueryParser(&req)

	reviews, err := h.Service.getFlaggedReviews(ctx, req.GameId, req.Limit, req.Offset)

	if err != nil {
		return GetFlaggedReviewsErrorResponse(c, err)
	}

	return GetFlaggedReviewsSuccessResp(c, reviews)
}

func (h *Handler) FlagReview(ctx context.Context, c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.Service.flagReview(ctx, id, true)

	if err != nil {
		return FlagReviewErrorResponse(c, err)
	}

	return FlagReviewSuccessResp(c)

}

type GetReview func(ctx context.Context, id GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error)

func getReview(ctx context.Context, c *fiber.Ctx, isForGame bool, getReview GetReview) error {
	var req GetReviewsForGame

	var gameId string
	var userId string

	err := c.QueryParser(&req)

	if err != nil {
		log.Println(err)
		return GetReviewErrorResponse(c, ErrBadRequest)
	}

	if isForGame {
		gameId = c.Params("gameId")
	} else {
		userId = c.Params("userId")
	}

	req.GameId = gameId
	req.UserId = userId

	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	if req.SortBy.Key == "" {
		req.SortBy.Key = "createdAt"
		req.SortBy.Asc = false
	}

	reviews, err := getReview(ctx, req)

	if err != nil {
		return GetReviewErrorResponse(c, err)
	}

	return GetReviewsSuccessResp(c, reviews)
}
