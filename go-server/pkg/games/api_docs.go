package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

// HandleAddGenre godoc
//
// @Security BearerAuth
//
//	@Summary		Adds a new game genre
//	@Description	Adds a new game genre, the slug is generated from the title, and it must be unique
//	@Tags			games
//	@ID				addGenre
//	@Accept			json
//	@Produce		json
//
//	@Param			addGenre	body		games.AddGenreRequest 	true			"addGenre request"
//
//	@Success		201				{object}	main.JSONResult{data=games.AddGenreRes}	"Success"
//	@Failure		409				{object}	main.JSONErrorRes					"Genre with the same slug already exists"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		401			{object}	main.JSONErrorRes											"Unauthorized
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/games/genres/add [post]
func HandleAddGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.AddGenre(ctx, c)
	}
}
