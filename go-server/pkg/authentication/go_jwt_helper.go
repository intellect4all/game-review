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
	log.Println("Validating JWT ")

	tokenString := string(jwtToken)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		log.Println("Parsing token")

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Println("Error parsing token")
			return nil, ErrInvalidJWT
		}

		log.Println("Token parsed")
		key, err := security.GetPrivateKey("key.pem")

		if err != nil {
			log.Println("Error getting private key")
			return nil, err
		}

		log.Println("Private key obtained")
		return key.Public(), nil
	})

	log.Println("Token parsed")

	if err != nil {
		log.Println("Error parsing token")
		return &JwtClaims{}, err
	}

	log.Println("Token parsed successfully", token.Claims)

	claims := token.Claims

	jwtClaims := &JwtClaims{}

	jwtClaims.fromMap(claims.(jwt.MapClaims))

	err = jwtClaims.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Claims parsed successfully", jwtClaims)

	if !token.Valid {
		log.Println("Error parsing claims")
		return &JwtClaims{}, err
	}

	log.Println("Claims parsed successfully")

	return jwtClaims, nil
}
