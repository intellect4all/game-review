package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type GameHandler struct {
	service *GameService
}

func NewGameHandler(service *GameService) *GameHandler {
	return &GameHandler{
		service: service,
	}
}

func (h *GameHandler) AddGenre(ctx context.Context, c *fiber.Ctx) error {
	var req AddGenreRequest

	err := c.BodyParser(&req)

	if err != nil {
		return AddGenreErrorResponse(c, ErrBadRequest)
	}

	gameGenre := &GameGenre{
		Title:     req.Title,
		Slug:      getSlug(req.Title),
		Desc:      req.Desc,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.service.AddGameGenre(ctx, gameGenre)

	if err != nil {
		return AddGenreErrorResponse(c, err)
	}

	return AddGenreSuccessResp(c, gameGenre.Slug)

}

type EditGenreRequest struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Desc  string `json:"desc"`
}

func (h *GameHandler) EditGenre(ctx context.Context, c *fiber.Ctx) error {
	var req EditGenreRequest

	err := c.BodyParser(&req)

	if err != nil {
		return EditGenreErrorResponse(c, ErrBadRequest)
	}

	gameGenre := &GameGenre{
		Title: req.Title,
		Slug:  req.Slug,
		Desc:  req.Desc,
	}

	err = h.service.EditGameGenre(ctx, gameGenre)

	if err != nil {
		return EditGenreErrorResponse(c, err)
	}

	return EditGenreSuccessResp(c)

}

func (h *GameHandler) GetGenres(ctx context.Context, c *fiber.Ctx) error {
	var req Pagination

	err := c.QueryParser(&req)

	if err != nil {
		return GetGenresErrorResponse(c, ErrBadRequest)
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	genres, err := h.service.GetAllGenres(ctx, &req)

	if err != nil {
		return GetGenresErrorResponse(c, err)
	}

	return GetGenresSuccessResp(c, genres)

}

func (h *GameHandler) GetGenre(ctx context.Context, c *fiber.Ctx) error {
	slug := c.Params("slug")

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return GetGenreErrorResponse(c, ErrGameGenreSlugRequired)
	}

	gameGenre, err := h.service.GetGameGenre(ctx, slug)

	if err != nil {
		return GetGenreErrorResponse(c, err)
	}

	return GetGenreSuccessResp(c, gameGenre)
}

func (h *GameHandler) DeleteGenre(ctx context.Context, c *fiber.Ctx) error {
	slug := c.Params("slug")

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return DeleteGenreErrorResponse(c, ErrGameGenreSlugRequired)
	}

	err := h.service.DeleteGameGenre(ctx, slug)

	if err != nil {
		return DeleteGenreErrorResponse(c, err)
	}

	return DeleteGenreSuccessResp(c)

}

func getSlug(title string) string {
	lowercase := strings.ToLower(title)

	return strings.ReplaceAll(lowercase, " ", "-")
}
