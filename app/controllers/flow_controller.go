package controllers

import (
	"InceptionAnimals/app/models"
	"InceptionAnimals/pkg/utils"
	"fmt"

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

func CreateNftTemplate(ctx *fiber.Ctx) error {

	template := models.NFTTemplate{}
	if err := ctx.BodyParser(&template); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_parse_json",
		})
	}

	validate := utils.NewValidator()
	// Validate user fields
	if err := validate.Struct(template); err != nil {
		// Return, if some fields are not valid.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	fmt.Printf("%+v\n", template)

	txRes, err := utils.CreateNftTemplate(template)
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
		"data":  txRes,
	})
}

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

	validate := utils.NewValidator()
	// Validate user fields
	if err := validate.Struct(body); err != nil {
		// Return, if some fields are not valid.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
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
