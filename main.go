package main

import (
	"InceptionAnimals/pkg/configs"
	"InceptionAnimals/pkg/middleware"
	"InceptionAnimals/pkg/routes"

	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()
	app := fiber.New(config)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Sanity Check!")
	})

	// Middlewares.
	middleware.FiberMiddleware(app)

	// Routes
	routes.UserRoutes(app) // Register a route for user.

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
