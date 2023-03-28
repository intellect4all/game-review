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
//	@Failure		404				{object}	main.JSONErrorRes											"Genre not found"
//	@Router			/api/v1/games/genres/update [put]
func HandleUpdateGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
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
//	@Success		200				{object}	main.JSONResult{data=games.PaginatedResponse[GameGenre]}	"Success"
//	@Router			/api/v1/games/genres [get]
func HandleGetGenres(handler *GameHandler, ctx context.Context) fiber.Handler {
	// set downstream context value
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
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
//	@Router			/api/v1/games/genres/{slug} [get]
func HandleGetGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
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
//	@Success		202				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"Genre not found"
//	@Router			/api/v1/games/genres/{slug} [delete]
func HandleDeleteGenre(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.DeleteGenre(ctx, c)
	}
}

// HandleAddGame godoc
//
// @Security BearerAuth
//
//	@Summary		Adds a new game
//	@Description	Adds a new game genre
//	@Tags			games
//	@ID				addGame
//	@Accept			json
//	@Produce		json
//
//	@Param			addGenre	body		games.AddGameRequest 	true			"addGame request"
//
//	@Success		201				{object}	main.JSONResult{data=games.AddGameRes}	"Success"
//	@Failure		409				{object}	main.JSONErrorRes					"Game already exists"
//	@Router			/api/v1/games/add [post]
func HandleAddGame(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.AddGame(ctx, c)
	}
}

// HandleGetGame godoc
//
// @Security BearerAuth
//
//	@Summary		Gets a game
//	@Description	the slug is required
//	@Tags			games
//	@ID				getGame
//	@Produce		json
//
//	@Param			id	path		string 	true			"id"
//
//	@Success		200				{object}	main.JSONResult{data=games.Game}	"Success"
//	@Router			/api/v1/games/{id} [get]
func HandleGetGame(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetGame(ctx, c)
	}
}

// HandleGetGames godoc
//
// @Security BearerAuth
//
//	@Summary		Gets all games
//	@Description	Gets all games, limits and offset can be used to paginate the results
//	@Tags			games
//	@ID				getGames
//	@Accept			json
//	@Produce		json
//
//	@Param			getGames	query		games.GetGamesQueries 	true			"getGames request"
//
//	@Success		200				{object}	main.JSONResult{data=games.PaginatedResponse[Game]}	"Success"
//	@Router			/api/v1/games [get]
func HandleGetGames(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.GetGames(ctx, c)
	}
}

// HandleUpdateGame godoc
//
// @Security BearerAuth
//
//	@Summary		Updates a game
//	@Description	the id is required
//	@Tags			games
//	@ID				updateGame
//	@Accept			json
//	@Produce		json
//
// @Param			id			path		string					true			"id"
//
//	@Param			updateGame	body		games.UpdateGameRequest 	true			"updateGame request"
//
//	@Success		200				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes					"Game not found"
//	@Router			/api/v1/games/{id} [put]
func HandleUpdateGame(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.UpdateGame(ctx, c)
	}
}

// HandleDeleteGame godoc
//
// @Security BearerAuth
//
//	@Summary		Delete a game
//	@Description	the id is required
//	@Tags			games
//	@ID				deleteGame
//	@Produce		json
//
//	@Param			id	path		string 	true			"id"
//
//	@Success		202				{object}	main.JSONResult{data=string}	"Success"
//	@Failure		404				{object}	main.JSONErrorRes											"Game not found"
//	@Router			/api/v1/games/{id} [delete]
func HandleDeleteGame(handler *GameHandler, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx = GetNewContext(ctx, c)
		return handler.DeleteGame(ctx, c)
	}
}

func GetNewContext(ctx context.Context, c *fiber.Ctx) context.Context {
	userId := c.Locals("userId").(string)
	role := c.Locals("role").(string)
	ctx = context.WithValue(ctx, "userId", userId)
	ctx = context.WithValue(ctx, "role", role)
	return ctx
}
