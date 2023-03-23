package games

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(mongoClient *mongo.Client, ctx context.Context, app fiber.Router, authNeeds *auth.AuthNeeds) error {

	gameRepo := NewGameRepositoryImpl(mongoClient)

	gameService := NewGameService(gameRepo)

	gameHandler := NewGameHandler(gameService)

	return router(ctx, app, gameHandler, authNeeds.AuthMiddleware)
}
