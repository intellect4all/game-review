package authentication

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type OTPCreationSuccessResponse struct {
	OtpID string `json:"otpId"`
	Email string `json:"email"`
}

func GetLoginSuccessResponse(d *LoginDTO) *fiber.Map {

	return getFiberMap(d, "Login successful")
}

func GetOTPCreationResponse(otp *string, email *string) *fiber.Map {
	resp := &OTPCreationSuccessResponse{
		OtpID: *otp,
		Email: *email,
	}

	return &fiber.Map{
		"message": "OTP sent successfully",
		"data":    resp,
	}
}

func getFiberMap(data *LoginDTO, message string) *fiber.Map {

	return &fiber.Map{
		"message": message,
		"data":    data,
	}
}

func SignUpErrorResponse(c *fiber.Ctx, err error) error {
	status := 0
	message := ""

	if err == ErrBadRequest {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	} else if err == ErrUserAlreadyExists {
		status = fiber.StatusConflict
		message = "User already exists"
	} else if err == ErrUsernameAlreadyExists {
		status = fiber.StatusConflict
		message = "Username already exists"
	} else {
		status = 500
		message = "Something went wrong"
	}
	log.Println(err)
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
	})
}
