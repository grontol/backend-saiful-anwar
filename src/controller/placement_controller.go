package controller

import (
	"yard_plan/src/response"
	"yard_plan/src/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PlacementController struct {
	placementService *service.PlacementService
	validate *validator.Validate
}

func NewPlacementController(
	placementService *service.PlacementService,
	validate *validator.Validate,
) *PlacementController {
	return &PlacementController{
		placementService: placementService,
		validate: validate,
	}
}

func (o *PlacementController) List(c *fiber.Ctx) error {	
	data, err := o.placementService.List()
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *PlacementController) ListByBlock(c *fiber.Ctx) error {
	id := c.Params("block_id", "")
	
	data, err := o.placementService.ListByBlock(id)
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}