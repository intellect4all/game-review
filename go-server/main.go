package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go-server/authentication"
	"os"
)

func init() {

}

func main() {
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
	ctx = context.WithValue(ctx, "apiVersion", "/v1")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello Here")
	})
	// register a ping route
	apiGroup.Get("/v1/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	err = authentication.Register(initResponse.mongoDbClient, ctx, apiGroup)
	if err != nil {
		return
	}

	err = app.Listen(":3000")
	if err != nil {
		os.Exit(1)
	}
}
