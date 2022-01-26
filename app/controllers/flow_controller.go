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
		})
	}
	// Return status 200 OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  latestBlock,
	})
}

// func CreateNftTemplate(ctx *fiber.Ctx) error {

// 	scriptName := "checkIsTemplateMinted"

// 	scriptRes, err := utils.ExecuteScript(scriptName)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": false,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	// Return status 200 OK
// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"error": false,
// 		"msg":   nil,
// 		"data":  scriptRes,
// 	})
// }

func CheckIfTemplateIsMinted(ctx *fiber.Ctx) error {
	type request struct {
		TypeID int `json:"typeId" validate:"required"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_parse_json",
		})
	}

	typeID := uint64(body.TypeID)
	scriptRes, err := utils.CheckIfTemplateIsMinted(typeID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": false,
			"msg":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  scriptRes,
	})
}
