package response

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Yard struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type Block struct {
	Id     int    `db:"id" json:"id"`
	YardId int    `db:"yard_id" json:"yard_id"`
	Name   string `db:"name" json:"name"`
	Slots  int    `db:"slots" json:"slots"`
	Rows   int    `db:"rows" json:"rows"`
	Tiers  int    `db:"tiers" json:"tiers"`
}

type YardPlan struct {
	Id           int     `db:"id" json:"id"`
	BlockId      int     `db:"block_id" json:"block_id"`
	SlotStart    int     `db:"slot_start" json:"slot_start"`
	SlotEnd      int     `db:"slot_end" json:"slot_end"`
	RowStart     int     `db:"row_start" json:"row_start"`
	RowEnd       int     `db:"row_end" json:"row_end"`
	Size         int     `db:"size" json:"size"`
	Height       float32 `db:"height" json:"height"`
	Type         string  `db:"type" json:"type"`
	SlotPriority int     `db:"slot_priority" json:"slot_priority"`
	RowPriority  int     `db:"row_priority" json:"row_priority"`
	TierPriority int     `db:"tier_priority" json:"tier_priority"`
}

type Placement struct {
	Id          int    `db:"id" json:"id"`
	ContainerId string `db:"container_id" json:"container_id"`
	BlockId     int    `db:"block_id" json:"block_id"`
	Slot        int    `db:"slot" json:"slot"`
	Row         int    `db:"row" json:"row"`
	Tier        int    `db:"tier" json:"tier"`
	IsHead      bool   `db:"is_head" json:"is_head"`
}

type Suggestion struct {
	BlockId int `db:"block_id" json:"block_id"`
	Slot    int `db:"slot" json:"slot"`
	Row     int `db:"row" json:"row"`
	Tier    int `db:"tier" json:"tier"`
}

func Success[T any](c *fiber.Ctx, data T) error {
	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"success": true,
		"message": "Success",
		"data":    data,
	})
}

func Error(c *fiber.Ctx, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"success": false,
			"message": validationErrors.Error(),
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(map[string]any{
			"success": false,
			"message": fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{
		"success": false,
		"message": fmt.Sprintf("Internal server error : %s", err),
	})
}
