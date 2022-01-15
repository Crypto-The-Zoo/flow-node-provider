package routes

import (
	"InceptionAnimals/app/controllers"

	"github.com/gofiber/fiber"
)

// UserRoutes func for describe group of user routes.
func UserRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method
	route.Get("/user/:id", controllers.GetUser)

	// Routes for POST method
	route.Post("/user/register", controllers.Register)
	route.Post("/user/login", controllers.Login)

	// Deprecated:
	// route.Post("/user/code", controllers.GetLoginCode)
	// // Routes for GET method
	// route.Get("/user/me", middleware.JWTProtected(), controllers.GetUser)

}
