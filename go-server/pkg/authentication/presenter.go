package authentication

import (
	"github.com/gofiber/fiber/v2"
)

type LoginSuccessResponse struct {
	User *UserDetail           `json:"user"`
	Jwt  *AuthenticatedUserJWT `json:"jwt"`
}

type OTPCreationSuccessResponse struct {
	OtpID string `json:"otpId"`
	Email string `json:"email"`
}

func GetLoginSuccessResponse(jwt *AuthenticatedUserJWT, detail *UserDetail) *fiber.Map {
	resp := &LoginSuccessResponse{
		Jwt:  jwt,
		User: detail,
	}
	return getFiberMap(resp, "Login successful")
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

func getFiberMap(data interface{}, message string) *fiber.Map {
	return &fiber.Map{
		"message": message,
		"data":    data,
	}
}
