package reviews

import (
	"context"
	"github.com/gofiber/fiber/v2"
	auth "go-server/pkg/authentication"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(mongoClient *mongo.Client, ctx context.Context, app fiber.Router, authNeeds *auth.AuthNeeds) error {

	repo := NewRepository(mongoClient)

	service := NewService(repo)

	handler := NewHandler(service)

	return router(ctx, app, handler, authNeeds.AuthMiddleware)
}
