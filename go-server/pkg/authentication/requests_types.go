package authentication

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignUpRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Username  string   `json:"username" validate:"required,ascii"`
	Password  string   `json:"password" validate:"required"`
	Role      string   `json:"role" validate:"required,oneof=user moderator"`
	FirstName string   `json:"firstName" validate:"required,ascii"`
	LastName  string   `json:"lastName" validate:"required,ascii"`
	Phone     string   `json:"phone" validate:"required,e164"`
	Location  Location `json:"location" validate:"required"`
}

type ForgetAndResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	TokenId         string `json:"tokenId" validate:"required"`
	OTPCode         string `json:"otpCode" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type VerifyAccountRequest struct {
	TokenID string `json:"tokenID" validate:"required"`
	OTPCode string `json:"otpCode" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}
