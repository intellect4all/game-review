package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"strings"
	"time"
)

type GameHandler struct {
	service *Service
}

func NewGameHandler(service *Service) *GameHandler {
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

func (h *GameHandler) AddGame(ctx context.Context, c *fiber.Ctx) error {
	var req AddGameRequest

	err := c.BodyParser(&req)

	if err != nil {
		return AddGameErrorResponse(c, ErrBadRequest)
	}

	game := &Game{
		Title:       req.Title,
		Summary:     req.Summary,
		ReleaseDate: req.ReleaseDate,
		Developer:   req.Developer,
		Publisher:   req.Publisher,
		Genres:      req.Genres,
		Rating:      RatingStats{},
		CreatedAt:   time.Now(),
	}

	err = h.service.AddGame(ctx, game)

	if err != nil {
		return AddGameErrorResponse(c, err)
	}

	return AddGameSuccessResp(c, game.Id.Hex())

}

func (h *GameHandler) GetGame(ctx context.Context, c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return GetGameErrorResponse(c, ErrGameIdRequired)
	}

	game, err := h.service.GetGame(ctx, id)

	if err != nil {
		return GetGameErrorResponse(c, err)
	}

	return GetGameSuccessResp(c, game)

}

type GetGamesQueries struct {
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	ReleasedDate string `json:"released_date,omitempty"`
	Developer    string `json:"developer,omitempty"`
	Publisher    string `json:"publisher,omitempty"`
	Genre        string `json:"genre,omitempty"`
}

func (h *GameHandler) GetGames(ctx context.Context, c *fiber.Ctx) error {
	log.Println("getting in query: ")
	var req GetGamesQueries

	err := c.QueryParser(&req)

	if err != nil {
		log.Println("getting in query error: ", err)
		return GetGamesErrorResponse(c, ErrBadRequest)
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	if req.Limit > 100 {
		req.Limit = 100
	}
	pagination := Pagination{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	log.Println("getting in query: ", req)

	filters := make(map[string]interface{})

	if req.ReleasedDate != "" {
		year, mon, day := 0, 0, 0
		rdStr := strings.TrimSpace(req.ReleasedDate)
		// break by "-"
		rdArr := strings.Split(rdStr, "-")

		for i, v := range rdArr {
			rdArr[i] = strings.TrimSpace(v)
			switch i {
			case 1:
				year, _ = strconv.Atoi(v)
			case 2:
				mon, _ = strconv.Atoi(v)
			case 3:
				day, _ = strconv.Atoi(v)
			}
		}

		releaseDate := time.Date(year, time.Month(mon), day, 0, 0, 0, 0, time.UTC)
		filters["released_date"] = releaseDate
	}

	log.Println("getting in query after rd: ", filters)

	if req.Developer != "" {
		filters["developer"] = req.Developer
	}
	log.Println("getting in query after d: ", filters)

	if req.Publisher != "" {
		filters["publisher"] = req.Publisher
	}
	log.Println("getting in query after p: ", filters)

	if req.Genre != "" {
		filters["genres.slug"] = req.Genre
	}

	log.Println("getting in query after g: ", filters)

	pagination.QueryFilters = filters

	games, err := h.service.GetAllGames(ctx, &pagination)

	if err != nil {
		return GetGamesErrorResponse(c, err)
	}

	return GetGamesSuccessResp(c, games)
}

func (h *GameHandler) UpdateGame(ctx context.Context, c *fiber.Ctx) error {
	idString := c.Params("id")

	if idString == "" {
		return UpdateGameErrorResp(c, ErrGameIdRequired)
	}

	var req UpdateGameRequest

	err := c.BodyParser(&req)

	if err != nil {
		return UpdateGameErrorResp(c, ErrBadRequest)
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return UpdateGameErrorResp(c, ErrBadRequest)
	}

	game := &Game{
		Id:          id,
		Title:       req.Title,
		Summary:     req.Summary,
		ReleaseDate: req.ReleaseDate,
		Developer:   req.Developer,
		Publisher:   req.Publisher,
		Genres:      req.Genres,
	}

	err = h.service.UpdateGame(ctx, game)

	if err != nil {
		return UpdateGameErrorResp(c, err)
	}

	return UpdateGameSuccessResp(c)

}

func (h *GameHandler) DeleteGame(ctx context.Context, c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return DeleteGameErrorResponse(c, ErrGameIdRequired)
	}

	err := h.service.DeleteGame(ctx, id)

	if err != nil {
		return DeleteGameErrorResponse(c, err)
	}

	return DeleteGameSuccessResp(c)
}

func getSlug(title string) string {
	lowercase := strings.ToLower(title)

	return strings.ReplaceAll(lowercase, " ", "-")
}
