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

		return nil, err
	}

	tokenString, err := token.SignedString(key)

	if err != nil {

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

	tokenString := string(jwtToken)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidJWT
		}

		key, err := security.GetPrivateKey("key.pem")

		if err != nil {
			return nil, err
		}

		return key.Public(), nil
	})

	if err != nil {
		return &JwtClaims{}, err
	}

	claims := token.Claims

	jwtClaims := &JwtClaims{}

	jwtClaims.fromMap(claims.(jwt.MapClaims))

	err = jwtClaims.Validate()
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return &JwtClaims{}, err
	}

	log.Println("JWt validated")

	return jwtClaims, nil
}
