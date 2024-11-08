package routes

import (
	"InceptionAnimals/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// FlowRoutes func for describe group of flow routes.
func FlowRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Post("/flow/block-events", controllers.GetEventsInBlockRange)
	route.Post("/flow/block-events-raw", controllers.GetEventsInBlockRangeRaw)
	// route.Post("/flow/template", middleware.JWTProtected(), controllers.CheckIfTemplateIsMinted)
	// route.Post("/flow/create-template", middleware.JWTProtected(), controllers.CreateNftTemplate)
	// route.Post("/flow/mint-nft", middleware.JWTProtected(), controllers.MintNFT)

}
