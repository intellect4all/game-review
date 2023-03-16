package authentication

import "github.com/gofiber/fiber/v2"

func OTPCreationSuccessResponse(otp, email string) *fiber.Map {
	return &fiber.Map{
		"message": "OTP sent successfully",
		"data":    map[string]string{"otpId": otp, "email": email},
	}
}
