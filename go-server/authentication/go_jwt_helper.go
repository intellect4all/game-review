package authentication

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"os"
)

type JWTHelperImpl struct{}

func NewJWTHelper() *JWTHelperImpl {
	return &JWTHelperImpl{}
}

func (j *JWTHelperImpl) GenerateJWT(claims *JwtClaims) (*AuthenticatedUserJWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	key, err := getJWTKey("key.pem")

	if err != nil {
		return nil, err
	}

	tokenString, err := token.SignedString(key)

	if err != nil {
		return nil, err
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

		return getJWTKey("key.pem.pub")
	})

	claims, ok := token.Claims.(JwtClaims)
	if !ok || !token.Valid {
		return &JwtClaims{}, err
	}

	return &claims, nil
}

func getJWTKey(file string) (*rsa.PrivateKey, error) {
	secretPrivateFile, err := os.Open(file)
	defer func(secretPrivateFile *os.File) {
		err := secretPrivateFile.Close()
		if err != nil {
			panic(err)
		}
	}(secretPrivateFile)

	if err != nil {
		return nil, err
	}

	fileStat, err := secretPrivateFile.Stat()

	if err != nil {
		return nil, err
	}

	secretKey := make([]byte, fileStat.Size())
	_, err = secretPrivateFile.Read(secretKey)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretKey)

	if err != nil {
		return nil, err
	}

	return privateKey, nil

}
