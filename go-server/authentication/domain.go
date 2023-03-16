package authentication

import (
	"context"
	"time"
)

type UserID string

type AuthenticatedUserJWT string

type User struct {
	Id   UserID `json:"email" bson:"email" validate:"required,email" `
	Role string `json:"role" bson:"role" validate:"required,oneof=customer admin"`
}

type UserCredential struct {
	Id         UserID    `json:"email" bson:"email" validate:"required,email" `
	Role       string    `json:"role" bson:"role" validate:"required,oneof=user admin moderator"`
	Password   string    `json:"password" bson:"password" validate:"required,min=8"`
	IsActive   bool      `json:"isActive" bson:"isActive"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	IsVerified bool      `json:"isVerified" bson:"isVerified"`
}

type UserDetail struct {
	User  UserCredential
	Phone string `json:"phone"bson:"phone" validate:"required"`
}

func (u *UserDetail) isAdmin() bool {
	return u.User.Role == "admin"
}

func (u *UserDetail) isCustomer() bool {
	return u.User.Role == "customer"
}

type Service struct {
	repository Repository
	jwtHelper  JWTHelper
}

func (s *Service) CreateUser(ctx context.Context, userCredential *UserCredential) error {
	return s.repository.CreateNewUser(ctx, userCredential)
}

func (s *Service) AuthenticateUser(ctx context.Context, userCredential *UserCredential) (*AuthenticatedUserJWT, error) {
	userCredentialFromDb, err := s.repository.GetUserCredential(ctx, userCredential.Id)
	if err != nil {
		return nil, err
	}

	if !s.repository.isValidPassword(userCredential.Password, userCredentialFromDb.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtHelper.GenerateJWT(userCredentialFromDb.Id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) refreshJWT(ctx context.Context, jwt AuthenticatedUserJWT) (token *AuthenticatedUserJWT, err error) {
	token, err = s.jwtHelper.RenewJWT(string(jwt))
	return
}

func (s *Service) DeleteUser(ctx context.Context, jwt AuthenticatedUserJWT, userID UserID) error {
	userDetail, err := s.jwtHelper.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	if !userDetail.isAdmin() {
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

	credential, err := s.GetUserCredential(ctx, *userId)

	tokenID = ""

	if err != nil {
		return "", err
	}

	if credential.IsVerified {
		return "", ErrUserAlreadyVerified
	}

	tokenID, err = s.repository.CreateVerificationOTP(ctx, userId)

	if err != nil {
		return "", err
	}

	return
}

type VerifyUserRequest struct {
	TokenID string `json:"tokenID" validate:"required"`
	OTPCode string `json:"otpCode" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

func (s *Service) VerifyUser(ctx context.Context, requestData *VerifyUserRequest) error {
	return s.repository.VerifyUser(ctx, requestData)

}

func NewService(repository Repository, jwtHelper JWTHelper) *Service {
	return &Service{repository: repository, jwtHelper: jwtHelper}
}

type Repository interface {
	CreateNewUser(ctx context.Context, userDetail *UserCredential) error
	GetUserCredential(ctx context.Context, email UserID) (*UserCredential, error)
	GetUserDetail(ctx context.Context, email UserID) (*UserDetail, error)
	DeleteUser(ctx context.Context, email UserID) error
	isValidPassword(passwordFromRequest, passwordFromDb string) bool
	CreateVerificationOTP(ctx context.Context, id *UserID) (string, error)
	VerifyUser(ctx context.Context, requestData *VerifyUserRequest) error
}

type JWTHelper interface {
	GenerateJWT(id UserID) (*AuthenticatedUserJWT, error)
	ValidateJWT(jwt AuthenticatedUserJWT) (UserDetail, error)
	RenewJWT(jwt string) (*AuthenticatedUserJWT, error)
}
