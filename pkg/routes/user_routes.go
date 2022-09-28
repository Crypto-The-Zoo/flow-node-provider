package routes

import (
	"github.com/gofiber/fiber/v2"
)

// UserRoutes func for describe group of user routes.
func UserRoutes(a *fiber.App) {
	// Create routes group.
	// route := a.Group("/api/v1")

	// // Routes for GET method
	// route.Get("/user/:id", controllers.GetUser)

	// // Routes for POST method
	// route.Post("/user/register", controllers.Register)
	// route.Post("/user/login", controllers.Login)

	// // Routes for PUT method
	// route.Put("/user/flow-address", middleware.JWTProtected(), controllers.AddFlowAddress)

	// Deprecated:
	// route.Post("/user/code", controllers.GetLoginCode)
	// // Routes for GET method
	// route.Get("/user/me", middleware.JWTProtected(), controllers.GetUser)

}
