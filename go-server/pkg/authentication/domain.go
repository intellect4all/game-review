package authentication

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
	"unicode"
)

type UserID string

type AuthenticatedUserJWT string

type JwtClaims struct {
	Email      string `json:"email"`
	Role       string `json:"role" validate:"required,oneof=user admin moderator"`
	IsVerified bool   `json:"isVerified"`
	jwt.StandardClaims
}

func (c *JwtClaims) Validate() error {
	return validator.New().Struct(c)
}

type UserCredential struct {
	Id         UserID      `json:"email" bson:"email" validate:"required,email" `
	Password   string      `json:"password" bson:"password" validate:"required,min=8"`
	IsActive   bool        `json:"isActive" bson:"isActive"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UserDetail *UserDetail `json:"userDetail" bson:"userDetail"`
}

type UserDetail struct {
	IsVerified bool   `json:"isVerified" bson:"isVerified"`
	FirstName  string `json:"firstName" bson:"firstName" validate:"required"`
	LastName   string `json:"lastName" bson:"lastName" validate:"required"`
	Phone      string `json:"phone" bson:"phone" validate:"omitempty"`
	Role       string `json:"role" bson:"role" validate:"required,oneof=user admin moderator"`
}

func NewService(repository Repository, jwtHelper JWTHelper) *Service {
	return &Service{repository: repository, jwtHelper: jwtHelper, validate: validator.New()}
}

type Repository interface {
	CreateNewUser(ctx context.Context, userDetail *UserCredential) error
	GetUserCredential(ctx context.Context, email UserID) (*UserCredential, error)
	GetUserDetail(ctx context.Context, email UserID) (*UserDetail, error)
	DeleteUser(ctx context.Context, email UserID) error
	CreateOTP(ctx context.Context, id *UserID) (string, error)
	VerifyUser(ctx context.Context, requestData *VerifyAccountRequest) error
	ChangePassword(ctx context.Context, f *ForgetAndResetPasswordRequest) error
}

type JWTHelper interface {
	GenerateJWT(claims *JwtClaims) (*AuthenticatedUserJWT, error)
	ValidateJWT(jwt AuthenticatedUserJWT) (*JwtClaims, error)
	RenewJWT(jwt AuthenticatedUserJWT) (*AuthenticatedUserJWT, error)
}

func (u *UserDetail) isAdmin() bool {
	return u.Role == "admin"
}

func (u *UserDetail) isCustomer() bool {
	return u.Role == "customer"
}

type Service struct {
	repository Repository
	jwtHelper  JWTHelper
	validate   *validator.Validate
}

func (s *Service) CreateUser(ctx context.Context, signUpRequest *SignUpRequest) error {
	valRes := s.validate.Struct(signUpRequest)

	if valRes != nil {
		return valRes
	}

	userID := UserID(strings.ToLower(signUpRequest.Email))

	// check if user already exists
	_, err := s.repository.GetUserCredential(ctx, userID)

	if err == nil {
		return ErrUserAlreadyExists
	}

	if ok, errMessage := isPasswordValid(signUpRequest.Password); !ok {
		return errors.New("Invalid Password: " + errMessage)
	}

	// encrypt password
	hashedPassword, err := encryptPassword(signUpRequest.Password)
	if err != nil {
		return UnknownError
	}

	userCred := getDefaultUserCredential(&userID, hashedPassword, *signUpRequest)

	err = s.repository.CreateNewUser(ctx, userCred)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AuthenticateUser(ctx context.Context, loginRequest *LoginRequest) (detail *UserDetail, jwtToken *AuthenticatedUserJWT, err error) {
	userId := UserID(strings.ToLower(loginRequest.Email))
	userCredentialFromDb, err := s.repository.GetUserCredential(ctx, userId)
	if err != nil {
		return
	}

	if !isCorrectPassword(loginRequest.Password, userCredentialFromDb.Password) {
		err = ErrInvalidCredentials
		return
	}

	if !userCredentialFromDb.IsActive {
		err = ErrAccountInactive
		return
	}

	detail = userCredentialFromDb.UserDetail

	claims := JwtClaims{
		Email:      string(userCredentialFromDb.Id),
		Role:       userCredentialFromDb.UserDetail.Role,
		IsVerified: userCredentialFromDb.UserDetail.IsVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth-service",
			Audience:  "game-reviews",
		},
	}

	log.Println("claims", claims)
	jwtToken, err = s.jwtHelper.GenerateJWT(&claims)
	return
}

func (s *Service) refreshJWT(jwt AuthenticatedUserJWT) (token *AuthenticatedUserJWT, err error) {
	token, err = s.jwtHelper.RenewJWT(jwt)
	return
}

func (s *Service) DeleteUser(ctx context.Context, jwt AuthenticatedUserJWT, userID UserID) error {
	claims, err := s.jwtHelper.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	if claims.Role != "admin" {
		return ErrUnauthorized
	}

	userToDelete, err := s.repository.GetUserDetail(ctx, userID)
	if err != nil {
		return err
	}

	if userToDelete.isAdmin() {
		return ErrUnauthorized
	}

	return s.repository.DeleteUser(ctx, userID)
}

func (s *Service) GetUserCredential(ctx context.Context, email UserID) (*UserCredential, error) {
	return s.repository.GetUserCredential(ctx, email)
}

func (s *Service) CreateVerificationOTP(ctx context.Context, userId *UserID) (tokenID string, err error) {
	if !s.isUserIDValid(userId) {
		return "", ErrInvalidRequest
	}
	credential, err := s.GetUserCredential(ctx, *userId)

	tokenID = ""

	if err != nil {
		return "", err
	}

	if credential.UserDetail.IsVerified {
		return "", ErrUserAlreadyVerified
	}

	tokenID, err = s.repository.CreateOTP(ctx, userId)

	if err != nil {
		return "", err
	}

	return
}

func (s *Service) VerifyUser(ctx context.Context, requestData *VerifyAccountRequest) error {
	return s.repository.VerifyUser(ctx, requestData)

}

func (s *Service) InitForgotPassword(ctx context.Context, userId *UserID) (string, error) {
	if !s.isUserIDValid(userId) {
		return "", ErrInvalidRequest
	}
	_, err := s.GetUserCredential(ctx, *userId)

	tokenID := ""

	if err != nil {
		return "", err
	}

	tokenID, err = s.repository.CreateOTP(ctx, userId)

	if err != nil {
		return "", err
	}

	return tokenID, nil
}

func (s *Service) ChangePassword(ctx context.Context, f *ForgetAndResetPasswordRequest) error {
	return s.repository.ChangePassword(ctx, f)
}

func isPasswordValid(password string) (bool, string) {
	var (
		upp, low, num, sym bool
		tot                uint8
		errorMessage       string
	)

	errorMessage = "Password must contain at least 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 symbol."

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			errorMessage = "Password contains invalid characters. Only letters, numbers, and symbols are allowed."
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {

		return false, errorMessage
	}

	return true, ""
}

func (s *Service) isUserIDValid(id *UserID) bool {
	err := s.validate.Var(id, "required,email")
	if err != nil {
		return false
	}
	return true
}

func getDefaultUserCredential(id *UserID, password string, request SignUpRequest) *UserCredential {

	return &UserCredential{
		Password:  password,
		IsActive:  strings.ToLower(request.Role) == "user",
		CreatedAt: time.Now(),
		Id:        *id,
		UserDetail: &UserDetail{
			Role:       trimAndLowercase(request.Role),
			IsVerified: false,
			FirstName:  strings.TrimSpace(request.FirstName),
			LastName:   strings.TrimSpace(request.LastName),
		},
	}

}

func trimAndLowercase(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isCorrectPassword(password string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
