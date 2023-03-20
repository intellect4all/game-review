package authentication

import (
	"github.com/golang-jwt/jwt"
	"go-server/pkg/security"
	"log"
)

type JWTHelperImpl struct{}

func NewJWTHelper() *JWTHelperImpl {
	return &JWTHelperImpl{}
}

func (j *JWTHelperImpl) GenerateJWT(claims *JwtClaims) (*AuthenticatedUserJWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	key, err := security.GetPrivateKey("key.pem")

	if err != nil {
		log.Printf("Error getting jwt key: %s", err.Error())
		return nil, err
	}

	tokenString, err := token.SignedString(key)

	if err != nil {
		log.Printf("Error generating jwt %s", err.Error())
		return nil, UnknownError
	}
	jwtToken := AuthenticatedUserJWT(tokenString)
	return &jwtToken, nil
}

func (j *JWTHelperImpl) RenewJWT(jwt AuthenticatedUserJWT) (token *AuthenticatedUserJWT, err error) {
	claims, err := j.ValidateJWT(jwt)

	if err != nil {
		return nil, err
	}

	token, err = j.GenerateJWT(claims)
	return token, err
}

func (j *JWTHelperImpl) ValidateJWT(jwtToken AuthenticatedUserJWT) (*JwtClaims, error) {
	token, err := jwt.Parse(string(jwtToken), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidJWT
		}

		return security.GetPrivateKey("key.pem.pub")
	})

	claims, ok := token.Claims.(JwtClaims)
	if !ok || !token.Valid {
		return &JwtClaims{}, err
	}

	return &claims, nil
}
