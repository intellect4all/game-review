package authentication

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"
)

type AuthenticatedUserJWT string

type JwtClaims struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role" validate:"required,oneof=user admin moderator"`
	IsVerified bool   `json:"isVerified"`
	jwt.StandardClaims
}

func (c *JwtClaims) fromMap(claims map[string]interface{}) *JwtClaims {
	c.Id = claims["id"].(string)
	c.Role = claims["role"].(string)
	c.Email = claims["email"].(string)
	c.IsVerified = claims["isVerified"].(bool)
	c.Issuer = claims["iss"].(string)
	c.Audience = claims["aud"].(string)
	c.ExpiresAt = int64(claims["exp"].(float64))
	c.IssuedAt = int64(claims["iat"].(float64))
	return c
}

func (c *JwtClaims) Validate() error {
	return validator.New().Struct(c)
}

type UserCredential struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"  validate:"required"`
	Password   string             `json:"password" bson:"password" validate:"required,min=8"`
	IsActive   bool               `json:"isActive" bson:"isActive"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UserDetail *UserDetail        `json:"userDetail" bson:"userDetail,inline"`
	Email      string             `json:"email" validate:"required,email"`
}

type UserDetail struct {
	IsVerified bool   `json:"isVerified" bson:"isVerified" validate:"required"`
	FirstName  string `json:"firstName" bson:"firstName" validate:"required"`
	LastName   string `json:"lastName" bson:"lastName" validate:"required"`
	Phone      string `json:"phone" bson:"phone" validate:"omitempty"`
	Username   string `json:"username" bson:"username" validate:"required"`
	Role       string `json:"role" bson:"role" validate:"required,oneof=user admin moderator"`
	displayPic string `json:"displayPic" bson:"displayPic"`
}

func NewService(repository AuthRepository, jwtHelper JWTHelper) *AuthService {
	return &AuthService{repository: repository, jwtHelper: jwtHelper, validate: validator.New()}
}

type AuthRepository interface {
	CreateNewUser(ctx context.Context, userCred UserCredential) error
	GetUserCredentialByEmail(ctx context.Context, email string) (*UserCredential, error)
	GetUserCredentialById(ctx context.Context, id string) (*UserCredential, error)
	GetUserCredentialByUserName(ctx context.Context, username string) (*UserCredential, error)
	GetUserDetail(ctx context.Context, email string) (*UserDetail, error)
	DeleteUser(ctx context.Context, email string) error
	CreateOTP(ctx context.Context, email string) (string, error)
	VerifyUser(ctx context.Context, req VerifyAccountRequest) error
	ChangePassword(ctx context.Context, req ForgetAndResetPasswordRequest) error
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

type AuthService struct {
	repository AuthRepository
	jwtHelper  JWTHelper
	validate   *validator.Validate
}

func (s *AuthService) CreateUser(ctx context.Context, signUpRequest SignUpRequest) error {
	valRes := s.validate.Struct(signUpRequest)

	if valRes != nil {
		return valRes
	}

	userID := strings.ToLower(signUpRequest.Email)

	// check if user already exists
	_, err := s.repository.GetUserCredentialByEmail(ctx, userID)

	if err == nil {
		return ErrUserAlreadyExists
	}

	// check if username already exists

	_, err = s.repository.GetUserCredentialByUserName(ctx, signUpRequest.Username)

	if err == nil {
		return ErrUsernameAlreadyExists
	}

	if ok, errMessage := isPasswordValid(signUpRequest.Password); !ok {
		return errors.New("Invalid Password: " + errMessage)
	}

	// encrypt password
	hashedPassword, err := encryptPassword(signUpRequest.Password)
	if err != nil {
		return UnknownError
	}

	userCred := getDefaultUserCredential(hashedPassword, signUpRequest)

	err = s.repository.CreateNewUser(ctx, *userCred)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) AuthenticateUser(ctx context.Context, loginRequest *LoginRequest) (detail *UserDetail, jwtToken *AuthenticatedUserJWT, err error) {
	userId := strings.ToLower(loginRequest.Email)
	userCredentialFromDb, err := s.repository.GetUserCredentialByEmail(ctx, userId)
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
		Id:         userCredentialFromDb.Id.Hex(),
		Role:       userCredentialFromDb.UserDetail.Role,
		Email:      userCredentialFromDb.Email,
		IsVerified: userCredentialFromDb.UserDetail.IsVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth-service",
			Audience:  "game-reviews",
		},
	}

	jwtToken, err = s.jwtHelper.GenerateJWT(&claims)
	return
}

func (s *AuthService) refreshJWT(jwt AuthenticatedUserJWT) (token *AuthenticatedUserJWT, err error) {
	token, err = s.jwtHelper.RenewJWT(jwt)
	return
}

func (s *AuthService) DeleteUser(ctx context.Context, jwt AuthenticatedUserJWT, email string) error {
	claims, err := s.jwtHelper.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	if claims.Role != "admin" {
		return ErrUnauthorized
	}

	userToDelete, err := s.repository.GetUserDetail(ctx, email)
	if err != nil {
		return err
	}

	if userToDelete.isAdmin() {
		return ErrUnauthorized
	}

	return s.repository.DeleteUser(ctx, email)
}

func (s *AuthService) GetUserCredential(ctx context.Context, email string) (*UserCredential, error) {
	return s.repository.GetUserCredentialByEmail(ctx, email)
}

func (s *AuthService) CreateVerificationOTP(ctx context.Context, email string) (tokenID string, err error) {
	if !s.isEmailValid(&email) {
		return "", ErrInvalidRequest
	}
	credential, err := s.GetUserCredential(ctx, email)

	tokenID = ""

	if err != nil {
		return "", err
	}

	if credential.UserDetail.IsVerified {
		return "", ErrUserAlreadyVerified
	}

	tokenID, err = s.repository.CreateOTP(ctx, email)

	if err != nil {
		return "", err
	}

	return
}

func (s *AuthService) VerifyUser(ctx context.Context, requestData VerifyAccountRequest) error {
	return s.repository.VerifyUser(ctx, requestData)

}

func (s *AuthService) InitForgotPassword(ctx context.Context, email string) (string, error) {
	if !s.isEmailValid(&email) {
		return "", ErrInvalidRequest
	}
	_, err := s.GetUserCredential(ctx, email)

	tokenID := ""

	if err != nil {
		return "", err
	}

	tokenID, err = s.repository.CreateOTP(ctx, email)

	if err != nil {
		return "", err
	}

	return tokenID, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, f ForgetAndResetPasswordRequest) error {
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

func (s *AuthService) isEmailValid(email *string) bool {
	err := s.validate.Var(email, "required,email")
	if err != nil {
		return false
	}
	return true
}

func getDefaultUserCredential(password string, request SignUpRequest) *UserCredential {

	return &UserCredential{
		Password:  password,
		IsActive:  strings.ToLower(request.Role) == "user",
		CreatedAt: time.Now(),
		Id:        primitive.ObjectID{},
		Email:     trimAndLowercase(request.Email),
		UserDetail: &UserDetail{
			Role:       trimAndLowercase(request.Role),
			IsVerified: false,
			FirstName:  strings.TrimSpace(request.FirstName),
			LastName:   strings.TrimSpace(request.LastName),
			Username:   strings.TrimSpace(request.Username),
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
