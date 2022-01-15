package main

import (
	"InceptionAnimals/pkg/configs"
	"InceptionAnimals/pkg/middleware"
	"InceptionAnimals/pkg/routes"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {

	// Load env variable
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env)
	fmt.Printf("--Env: %s", env)

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

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}
