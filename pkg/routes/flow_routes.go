package routes

import (
	"InceptionAnimals/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// FlowRoutes func for describe group of flow routes.
func FlowRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Get("/flow/block", controllers.GetLatestBlock)

	route.Post("/flow/template", controllers.CheckIfTemplateIsMinted)

	// route.Post("/flow/nft-template", controllers.CreateNftTemplate)

}
