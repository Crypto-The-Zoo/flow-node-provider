package main

import (
	"InceptionAnimals/pkg/routes"

	"github.com/gofiber/fiber"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Sanity Check!")
	})
	// middleware.FiberMiddleware(app)

	// Routes
	routes.UserRoutes(app) // Register a route for user.

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}
}
