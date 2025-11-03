package controller

import (
	"yard_plan/src/response"
	"yard_plan/src/service"
	"yard_plan/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type YardController struct {
	yardService *service.YardService
	validate *validator.Validate
}

func NewYardController(
	yardService *service.YardService,
	validate *validator.Validate,
) *YardController {
	return &YardController{
		yardService: yardService,
		validate: validate,
	}
}

func (o *YardController) List(c *fiber.Ctx) error {	
	data, err := o.yardService.List()
	if err != nil {
		return err
	}
	
	return response.Success(c, data)
}

func (o *YardController) Create(c *fiber.Ctx) error {
	req := new(validation.YardPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	if err := o.yardService.Create(req); err != nil {
		return err
	}
	
	return response.Success(c, "Inserted")
}

func (o *YardController) Edit(c *fiber.Ctx) error {
	id := c.Params("id", "")
	
	req := new(validation.YardPayload)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := o.validate.Struct(req); err != nil {
		return err
	}
	
	if err := o.yardService.Edit(id, req); err != nil {
		return err
	}
	
	return response.Success(c, "Updated")
}

func (o *YardController) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	
	if err := o.yardService.Delete(id); err != nil {
		return err
	}
	
	return response.Success(c, "Deleted")
}