package reviews

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
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
	Rating   int      `json:"rating" validate:"required,number,gte=0,lte2"`
	Comment  string   `json:"comment" validate:"required,min=5,max=2000"`
	GameId   string   `json:"gameId" validate:"required"`
	Location Location `json:"location" validate:"required"`
}

func (h *Handler) AddReview(ctx context.Context, c *fiber.Ctx) error {
	var req AddReviewRequest

	err := c.BodyParser(&req)

	if err != nil {
		return AddReviewErrorResponse(c, ErrBadRequest)
	}
	userId := c.Locals("userId").(string)
	review := &AddReview{
		Rating:   req.Rating,
		Comment:  req.Comment,
		GameId:   req.GameId,
		UserId:   userId,
		Location: req.Location,
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
	GameId string `json:"gameId"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (h *Handler) GetFlaggedReviews(ctx context.Context, c *fiber.Ctx) error {
	var req GetFlaggedReviewsRequest

	log.Println("here domain 0")

	fmt.Println("here domain 0")

	err := c.QueryParser(&req)

	log.Println("here domain 1")

	reviews, err := h.Service.getFlaggedReviews(ctx, req.GameId, req.Limit, req.Offset)

	log.Println("here domain 2")

	if err != nil {
		log.Println(err)
		return GetFlaggedReviewsErrorResponse(c, err)
	}

	return GetFlaggedReviewsSuccessResp(c, reviews)
}

func (h *Handler) FlagReview(ctx context.Context, c *fiber.Ctx, shouldFlag bool) error {
	id := c.Params("id")

	err := h.Service.flagReview(ctx, id, shouldFlag)

	if err != nil {
		return FlagReviewErrorResponse(c, err)
	}

	return FlagReviewSuccessResp(c, shouldFlag)

}

type GetReviewsLocationsRequest struct {
	Type  string `json:"type" bson:"type" validate:"oneof=day week month"`
	Value int    `json:"value" bson:"value" validate:"required,number,gte=1"`
}

func (h *Handler) GetReviewsLocations(ctx context.Context, c *fiber.Ctx) error {
	var req GetReviewsLocationsRequest

	err := c.QueryParser(&req)

	if err != nil {
		return GetLocationErrorResponse(c, ErrBadRequest)
	}

	if validator.New().Struct(req) != nil {
		return GetLocationErrorResponse(c, ErrBadRequest)
	}

	var reqType LocationReqType

	if req.Type == "day" {
		reqType = Day
	} else if req.Type == "week" {
		reqType = Week
	} else if req.Type == "month" {
		reqType = Month
	} else {
		reqType = Day
	}

	locations, err := h.Service.getLocation(ctx, reqType, req.Value)

	if err != nil {
		return GetLocationErrorResponse(c, err)
	}

	return GetLocationSuccessResp(c, locations)

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
