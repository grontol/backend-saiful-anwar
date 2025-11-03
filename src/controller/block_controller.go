package controller

import (
	"yard_plan/src/response"
	"yard_plan/src/service"
	"yard_plan/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BlockController struct {
	blockService *service.BlockService
	validate *validator.Validate
}

func NewBlockController(
	blockService *service.BlockService,
	validate *validator.Validate,
) *BlockController {
	return &BlockController{
		blockService: blockService,
		validate: validate,
	}
}

func (o *BlockController) List(c *fiber.Ctx) error {	
	data, err := o.blockService.List()
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *BlockController) ListByYard(c *fiber.Ctx) error {
	id := c.Params("yard_id", "")
	
	data, err := o.blockService.ListByYard(id)
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *BlockController) Create(c *fiber.Ctx) error {
	req := new(validation.BlockPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	if err := o.blockService.Create(req); err != nil {
		return err
	}
	
	return response.Success(c, "Inserted")
}

func (o *BlockController) Edit(c *fiber.Ctx) error {
	id := c.Params("id", "")
	
	req := new(validation.BlockPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	if err := o.blockService.Edit(id, req); err != nil {
		return err
	}
	
	return response.Success(c, "Updated")
}

func (o *BlockController) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	
	if err := o.blockService.Delete(id); err != nil {
		return err
	}
	
	return response.Success(c, "Deleted")
}