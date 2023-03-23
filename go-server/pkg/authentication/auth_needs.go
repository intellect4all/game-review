package authentication

type AuthNeeds struct {
	JwtHelper      *JWTHelperImpl
	AuthChecker    *AuthChecker
	AuthMiddleware *AuthMiddlewareImpl
}

func NewAuthNeeds() *AuthNeeds {
	return &AuthNeeds{
		JwtHelper:      NewJWTHelper(),
		AuthChecker:    NewAuthChecker(),
		AuthMiddleware: NewAuthMiddleware(NewJWTHelper()),
	}
}
