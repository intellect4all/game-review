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

func GetGenresSuccessResp(c *fiber.Ctx, genreResp *PaginatedGameGenres) error {
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
