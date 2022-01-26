package controllers

import (
	"InceptionAnimals/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func GetLatestBlock(ctx *fiber.Ctx) error {

	latestBlock, err := utils.GetLatestBlock()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": false,
			"msg":   err.Error(),
			"data":  "unable to get latest block from Flow",
		})
	}
	// Return status 200 OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  latestBlock,
	})
}
