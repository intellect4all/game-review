package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberAdapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go-server/authentication"
	"os"
)

var fiberLambda *fiberAdapter.FiberLambda

func init() {
	// get the environment variables and initialize the application
	initResponse, err := InitializationHandler()
	if err != nil {
		panic(err)
	}

	var app *fiber.App
	app = fiber.New()

	// log the requests
	app.Use(logger.New())
	app.Use(recover.New())

	apiGroup := app.Group("/api/")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "apiVersion", "v1")

	// register a ping route
	apiGroup.Get("v1/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	authentication.Register(initResponse.mongoDbClient, ctx, apiGroup)

	if !initResponse.environmentIsLocal {
		fiberLambda = fiberAdapter.New(app)
		return
	}

	err = app.Listen(":3000")
	if err != nil {
		os.Exit(1)
	}
}

func main() {

	if fiberLambda != nil {
		lambda.Start(Handler)
	}

}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}
