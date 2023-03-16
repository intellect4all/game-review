package authentication

import "errors"

var ErrInvalidCredentials = errors.New("invalid-credentials")
var ErrUnauthorized = errors.New("unauthorized")
var ErrUserNotFound = errors.New("user-not-found")
var ErrUserAlreadyExists = errors.New("user-already-exists")
var ErrUserAlreadyVerified = errors.New("user-already-verified")
var ErrOTPCreationFailed = errors.New("token-creation-failed")
var ErrInvalidOTP = errors.New("invalid-otp")
var ErrOTPUsed = errors.New("otp-used")
var ErrOTPExpired = errors.New("otp-expired")
var UnknownError = errors.New("unknown-error")
var ErrInvalidRequest = errors.New("invalid-request")