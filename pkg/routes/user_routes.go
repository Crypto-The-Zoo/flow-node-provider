package routes

import (
	"InceptionAnimals/app/controllers"
	"InceptionAnimals/pkg/middleware"

	"github.com/gofiber/fiber"
)

// UserRoutes func for describe group of user routes.
func UserRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method
	route.Get("/user/:id", middleware.JWTProtected(), controllers.GetUser)
	route.Get("/user/me", middleware.JWTProtected(), controllers.GetUser)

	// Routes for POST method
	route.Post("/user/register", controllers.CreateUser)
	route.Post("/user/login", controllers.Login)
}
