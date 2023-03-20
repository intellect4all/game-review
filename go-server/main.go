package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "go-server/docs"
	"go-server/pkg/authentication"
	"log"
	"os"
)

//	@title			Game Review API
//	@version		1.0
//	@description	This is an Api Service for Cool Game Review Api.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	autolarry55@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:3000
// @BasePath	/
// @schemes	http
func main() {
	// get the environment variables and initialize the application
	initResponse, err := InitializationHandler()
	if err != nil {
		panic(err)
	}

	var app *fiber.App
	app = fiber.New()

	// log the requests
	app.Use(cors.New())
	//app.Use(csrf.New())
	app.Use(logger.New())
	app.Use(recover.New())

	apiGroup := app.Group("/api/")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "apiVersion", "/v1")

	fmt.Println("Server is running on port 3000")
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/loaderio-baf5658b393ffe75ded7e5209eb81d79.txt", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("resources/loaderTest.txt")

	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello Here")
	})
	// register a ping route
	apiGroup.Get("/v1/ping", ping)

	err = authentication.Register(initResponse.MongoDbClient, ctx, apiGroup)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = app.Listen(":3000")
	if err != nil {
		os.Exit(1)
	}
}

// Ping godoc
//
//	@Summary		Show the status of server.
//	@Description	get the status of server.
//
//	@ID				ping
//
//	@Tags			ping
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/api/v1/ping [get]
func ping(ctx *fiber.Ctx) error {
	res := map[string]interface{}{
		"status":  "success",
		"result":  "pong",
		"message": "Server is up and running",
	}

	if err := ctx.JSON(res); err != nil {
		return err
	}

	return nil

}

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONErrorRes struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}
