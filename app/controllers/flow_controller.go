package controllers

import (
	"InceptionAnimals/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func GetEventsInBlockRangeRaw(ctx *fiber.Ctx) error {
	type request struct {
		Node        string `json:"node" validate:"required"`
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
	node := body.Node
	startHeight := uint64(body.StartHeight)
	endHeight := uint64(body.EndHeight)

	blocks, err := utils.GetEventsInBlockHeightRangeRaw(node, eventType, startHeight, endHeight)
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