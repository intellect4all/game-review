package games

import "github.com/gofiber/fiber/v2"

func AddGenreErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameGenreAlreadyExists {
		status = fiber.StatusConflict
		message = "Game genre already existed"
	} else if err == ErrUnauthorized {
		status = fiber.StatusUnauthorized
		message = "Unauthorized"
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
