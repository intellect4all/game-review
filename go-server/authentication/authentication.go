package authentication

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(mongoClient *mongo.Client, ctx context.Context, app fiber.Router) error {
	authRepo := NewMongoRepository(mongoClient)
	jwtHelper := NewJWTHelper()
	authService := NewService(authRepo, jwtHelper)
	authHandler := NewHandler(authService)

	return authRouter(ctx, app, authHandler)
}

//func BookSuccessResponse(data *entities.Book) *fiber.Map {
//	book := Book{
//		ID:     data.ID,
//		Title:  data.Title,
//		Author: data.Author,
//	}
//	return &fiber.Map{
//		"status": true,
//		"data":   book,
//		"error":  nil,
//	}
//}
//
//// BooksSuccessResponse is the list SuccessResponse that will be passed in the response by Handler
//func BooksSuccessResponse(data *[]Book) *fiber.Map {
//	return &fiber.Map{
//		"status": true,
//		"data":   data,
//		"error":  nil,
//	}
//}
//
//// BookErrorResponse is the ErrorResponse that will be passed in the response by Handler
//func BookErrorResponse(err error) *fiber.Map {
//	return &fiber.Map{
//		"status": false,
//		"data":   "",
//		"error":  err.Error(),
//	}
//}
