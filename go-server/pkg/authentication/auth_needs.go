package authentication

type AuthNeeds struct {
	JwtHelper      *JWTHelperImpl
	AuthMiddleware *AuthMiddlewareImpl
}

func NewAuthNeeds() *AuthNeeds {
	return &AuthNeeds{
		JwtHelper:      NewJWTHelper(),
		AuthMiddleware: NewAuthMiddleware(NewJWTHelper()),
	}
}
