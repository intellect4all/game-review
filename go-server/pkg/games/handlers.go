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
		DateAdded: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.service.AddGameGenre(ctx, gameGenre)

	if err != nil {
		return AddGenreErrorResponse(c, err)
	}

	return AddGenreSuccessResp(c, gameGenre.Slug)

}

func getSlug(title string) string {
	lowercase := strings.ToLower(title)

	return strings.ReplaceAll(lowercase, " ", "-")

}
