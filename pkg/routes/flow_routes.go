package routes

import (
	"InceptionAnimals/app/controllers"
	"InceptionAnimals/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// FlowRoutes func for describe group of flow routes.
func FlowRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Get("/flow/block", controllers.GetLatestBlock)

	route.Post("/flow/template", middleware.JWTProtected(), controllers.CheckIfTemplateIsMinted)
	route.Post("/flow/create-template", middleware.JWTProtected(), controllers.CreateNftTemplate)

}
