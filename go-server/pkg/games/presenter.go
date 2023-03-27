package games

import (
	"github.com/gofiber/fiber/v2"
)

func AddGenreErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameGenreAlreadyExists {
		status = fiber.StatusConflict
		message = "Game genre already existed"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

type AddGenreRes struct {
	Slug string `json:"slug"`
}

func AddGenreSuccessResp(c *fiber.Ctx, slug string) error {
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Game genre added",
		"data":    AddGenreRes{Slug: slug},
	})

}

func EditGenreErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameGenreSlugRequired {
		status = fiber.StatusBadRequest
		message = "Game genre slug is required"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game genre not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func EditGenreSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Game genre updated",
		"data":    "",
	})

}

func GetGenresErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func GetGenresSuccessResp(c *fiber.Ctx, genreResp *PaginatedResponse[GameGenre]) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Game genres",
		"data":    genreResp,
	})
}

func GetGenreErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameGenreSlugRequired {
		status = fiber.StatusBadRequest
		message = "Game genre slug is required"

	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game genre not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func GetGenreSuccessResp(c *fiber.Ctx, genre *GameGenre) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Game genre",
		"data":    genre,
	})
}

func DeleteGenreErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameGenreSlugRequired {
		status = fiber.StatusBadRequest
		message = "Game genre slug is required"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game genre not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func DeleteGenreSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Game genre deleted",
		"data":    "",
	})
}

func AddGameErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameAlreadyExists {
		status = fiber.StatusConflict
		message = "Game already existed"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

type AddGameRes struct {
	gameId string `json:"gameId"`
}

func AddGameSuccessResp(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Game added",
		"data":    AddGameRes{gameId: id},
	})
}

func GetGameErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameIdRequired {
		status = fiber.StatusBadRequest
		message = "Game id is required"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func GetGameSuccessResp(c *fiber.Ctx, game *Game) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Game",
		"data":    game,
	})
}

func GetGamesErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func GetGamesSuccessResp(c *fiber.Ctx, gamesResponse *PaginatedResponse[Game]) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Games",
		"data":    gamesResponse,
	})
}

func UpdateGameErrorResp(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameIdRequired {
		status = fiber.StatusBadRequest
		message = "Game id is required"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func UpdateGameSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Game updated",
		"data":    "",
	})
}

func DeleteGameErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameIdRequired {
		status = fiber.StatusBadRequest
		message = "Game id is required"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Game not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func DeleteGameSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Game deleted",
		"data":    "",
	})
}
