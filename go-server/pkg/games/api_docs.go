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

// HandleUpdateGenre godoc
//
// @Security BearerAuth
//
//	@Summary		Updates a game genre
//	@Description	Updates a game genre, the slug is required
//	@Tags			games
//	@ID				updateGenre
//	@Accept			json
//	@Produce		json
//
//	@Param			updateGenre	body		games.EditGenreRequest 	true			"updateGenre request"
//
//	@Success		202				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		401			{object}	main.JSONErrorRes											"Unauthorized
//	@Failure		404				{object}	main.JSONErrorRes											"Genre not found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/games/genres/update [post]
func HandleUpdateGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.EditGenre(ctx, c)
	}
}

// HandleGetGenres godoc
//
// @Security BearerAuth
//
//	@Summary		Gets all game genres
//	@Description	Gets all game genres, limits and offset can be used to paginate the results
//	@Tags			games
//	@ID				getGenres
//	@Accept			json
//	@Produce		json
//
//	@Param			getGenres	query		games.Pagination 	true			"getGenres request"
//
//	@Success		200				{object}	main.JSONResult{data=games.PaginatedGameGenres}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		401			{object}	main.JSONErrorRes											"Unauthorized
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/games/genres [get]
func HandleGetGenres(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.GetGenres(ctx, c)
	}
}

// HandleGetGenre godoc
//
// @Security BearerAuth
//
//	@Summary		Gets a game genre
//	@Description	the slug is required
//	@Tags			games
//	@ID				getGenre
//	@Produce		json
//
//	@Param			slug	path		string 	true			"slug"
//
//	@Success		200				{object}	main.JSONResult{data=games.GameGenre}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		401			{object}	main.JSONErrorRes											"Unauthorized
//	@Failure		404				{object}	main.JSONErrorRes											"Genre not found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/games/genres/{slug} [get]
func HandleGetGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.GetGenre(ctx, c)
	}
}

// HandleDeleteGenre godoc
//
// @Security BearerAuth
//
//	@Summary		Delete a game genre
//	@Description	the slug is required
//	@Tags			games
//	@ID				deleteGenre
//	@Produce		json
//
//	@Param			slug	path		string 	true			"slug"
//
//	@Success		200				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		400				{object}	main.JSONErrorRes											"Bad request"
//	@Failure		401			{object}	main.JSONErrorRes											"Unauthorized
//	@Failure		404				{object}	main.JSONErrorRes											"Genre not found"
//	@Failure		500				{object}	main.JSONErrorRes											"Internal Server Error"
//	@Router			/api/v1/games/genres/{slug} [delete]
func HandleDeleteGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler.DeleteGenre(ctx, c)
	}
}
