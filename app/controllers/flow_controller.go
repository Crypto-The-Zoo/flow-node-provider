package controllers

import (
	"InceptionAnimals/app/models"
	"InceptionAnimals/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetEventsInBlockRange(ctx *fiber.Ctx) error {
	type request struct {
		Type        string `json:"type" validate:"required"`
		StartHeight int    `json:"startHeight" validate:"required"`
		EndHeight   int    `json:"endHeight" validate:"required"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_parse_json",
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	eventType := body.Type
	startHeight := uint64(body.StartHeight)
	endHeight := uint64(body.EndHeight)

	blocks, err := utils.GetEventsInBlockHeightRangeAutoNode(eventType, startHeight, endHeight)
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
		"data":  blocks,
	})
}

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
	if err := validate.Struct(template); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}
	fmt.Printf("--Template: %+v\n", template)

	// Get claims from JWT
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		// Return status 500 and JWT parse error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if claims.Email != "alliu930410@gmail.com" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "admin_required",
		})
	}

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
	if err := validate.Struct(body); err != nil {
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

func MintNFT(ctx *fiber.Ctx) error {
	var body models.MintNFTRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_parse_json",
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get claims from JWT
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		// Return status 500 and JWT parse error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if claims.Email != "alliu930410@gmail.com" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "admin_required",
		})
	}

	txRes, err := utils.MintNFT(body.TypeID, body.Address)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data":  txRes,
	})

}
