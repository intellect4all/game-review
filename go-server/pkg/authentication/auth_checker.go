package authentication

import (
	"context"
	"log"
)

type AuthChecker struct {
}

func NewAuthChecker() *AuthChecker {
	return &AuthChecker{}
}

func (a *AuthChecker) IsAuthenticated(ctx context.Context) bool {
	claims := ctx.Value("userClaims")

	log.Println("claims: from is authenticated ", claims)
	if claims == nil {
		return false
	}

	return true
}

func (a *AuthChecker) IsAdmin(ctx context.Context) bool {
	isLogged := a.IsAuthenticated(ctx)

	if !isLogged {
		return false
	}

	claims := ctx.Value("userClaims").(JwtClaims)

	if claims.Role == "admin" || claims.Role == "moderator" {
		return true
	}

	return false
}
