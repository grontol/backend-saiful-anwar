package controller

import (
	"yard_plan/src/response"
	"yard_plan/src/service"
	"yard_plan/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type YardPlanController struct {
	yardPlanService *service.YardPlanService
	validate *validator.Validate
}

func NewYardPlanController(
	yardPlanService *service.YardPlanService,
	validate *validator.Validate,
) *YardPlanController {
	return &YardPlanController{
		yardPlanService: yardPlanService,
		validate: validate,
	}
}

func (o *YardPlanController) List(c *fiber.Ctx) error {	
	data, err := o.yardPlanService.List()
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *YardPlanController) ListByYard(c *fiber.Ctx) error {
	id := c.Params("yard_id", "")
	
	data, err := o.yardPlanService.ListByYard(id)
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *YardPlanController) ListByBlock(c *fiber.Ctx) error {
	id := c.Params("block_id", "")
	
	data, err := o.yardPlanService.ListByBlock(id)
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *YardPlanController) Create(c *fiber.Ctx) error {
	req := new(validation.YardPlanPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	if err := o.yardPlanService.Create(req); err != nil {
		return err
	}
	
	return response.Success(c, "Inserted")
}

func (o *YardPlanController) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	
	if err := o.yardPlanService.Delete(id); err != nil {
		return err
	}
	
	return response.Success(c, "Deleted")
}

func (o *YardPlanController) Suggest(c *fiber.Ctx) error {
	req := new(validation.SuggestionPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	suggestion, err := o.yardPlanService.Suggest(req)
	if err != nil {
		return err
	}
	
	if suggestion == nil {
		return fiber.NewError(fiber.StatusNotFound, "No available suggestion")
	}
	
	return response.Success(c, suggestion)
}

func (o *YardPlanController) Place(c *fiber.Ctx) error {
	req := new(validation.PlacementPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	err := o.yardPlanService.Place(req)
	if err != nil {
		return err
	}
	
	return response.Success(c, "Placemenet success")
}

func (o *YardPlanController) Pickup(c *fiber.Ctx) error {
	req := new(validation.PickupPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	err := o.yardPlanService.Pickup(req)
	if err != nil {
		return err
	}
	
	return response.Success(c, "Pickup success")
}
