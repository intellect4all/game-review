package reviews

import "github.com/gofiber/fiber/v2"

func AddReviewErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrGameNotFound {
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

type AddReviewRes struct {
	ReviewId string `json:"reviewId"`
}

func AddReviewSuccessResp(c *fiber.Ctx, reviewId string) error {
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Review added",
		"data":    AddReviewRes{ReviewId: reviewId},
	})
}

func UpdateReviewErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Review not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func UpdateReviewSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Review updated",
		"data":    "",
	})
}

func GetReviewErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Review not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func GetReviewSuccessResp(c *fiber.Ctx, review ReviewResponse) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Review found",
		"data":    review,
	})
}

func GetReviewsSuccessResp(c *fiber.Ctx, reviews *PaginatedResponse[ReviewResponse]) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Reviews found",
		"data":    reviews,
	})
}

func DeleteReviewSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Review deleted",
		"data":    "",
	})
}

func VoteReviewErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Review not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func VoteReviewSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Review voted",
		"data":    "",
	})
}

func GetFlaggedReviewsErrorResponse(c *fiber.Ctx, err error) error {
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

func GetFlaggedReviewsSuccessResp(c *fiber.Ctx, reviews *PaginatedResponse[Review]) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Flagged reviews response",
		"data":    reviews,
	})
}

func FlagReviewErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrNotFound {
		status = fiber.StatusNotFound
		message = "Review not found"
	} else {
		status = 500
		message = "Something went wrong"
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}

func FlagReviewSuccessResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"message": "Review flagged",
		"data":    "",
	})
}
