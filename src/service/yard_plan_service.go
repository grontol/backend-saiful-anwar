package service

import (
	"yard_plan/src/response"
	"yard_plan/src/utils"
	"yard_plan/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type YardPlanService struct {
	db               *sqlx.DB
	blockService     *BlockService
	placementService *PlacementService
}

func NewYardPlanService(
	db *sqlx.DB,
	blockService *BlockService,
	placementService *PlacementService,
) *YardPlanService {
	return &YardPlanService{
		db:               db,
		blockService:     blockService,
		placementService: placementService,
	}
}

func (o *YardPlanService) List() ([]response.YardPlan, error) {
	data := []response.YardPlan{}
	err := o.db.Select(&data, "SELECT * FROM yard_plans ORDER BY id")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *YardPlanService) ListByYard(yardId string) ([]response.YardPlan, error) {
	data := []response.YardPlan{}
	err := o.db.Select(
		&data,
		`
			SELECT yard_plans.* FROM yard_plans
			JOIN blocks ON blocks.id = yard_plans.block_id
			WHERE yard_id = $1 ORDER BY yard_plans.id
		`,
		yardId,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *YardPlanService) ListByBlock(blockId string) ([]response.YardPlan, error) {
	data := []response.YardPlan{}
	err := o.db.Select(&data, "SELECT * FROM yard_plans WHERE block_id = $1 ORDER BY id", blockId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *YardPlanService) Create(data *validation.YardPlanPayload) error {
	block, err := o.blockService.GetById(data.BlockId)

	if block == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Block not found")
	}

	if data.SlotStart < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Slot start must be greater than zero")
	}

	if data.SlotEnd > block.Slots {
		return fiber.NewError(fiber.StatusBadRequest, "Row start must be less than block slots")
	}

	if data.RowStart < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Row start must be greater than zero")
	}

	if data.RowEnd > block.Rows {
		return fiber.NewError(fiber.StatusBadRequest, "Row start must be less than block rows")
	}

	slotRequired := 1
	if data.Size == 40 {
		slotRequired = 2
	}

	if data.SlotEnd-data.SlotStart+1 < slotRequired {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid slot range")
	}

	if data.RowEnd-data.RowStart+1 < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid row range")
	}

	// Check for intersecting plan
	var intersectingPlans int
	o.db.Get(
		&intersectingPlans,
		`
			SELECT COUNT(*) FROM yard_plans
			WHERE NOT (slot_start > $1 OR slot_end < $2 OR row_start > $3 OR row_end < $4)
		`,
		data.SlotEnd,
		data.SlotStart,
		data.RowEnd,
		data.RowStart,
	)

	// Check availability
	var alreadyPlaced int
	o.db.Get(
		&alreadyPlaced,
		`
			SELECT COUNT(*) FROM placements
			WHERE slot >= $1 AND slot <= $2 AND row >= $3 AND row <= $4
		`,
		data.SlotStart,
		data.SlotEnd,
		data.RowStart,
		data.RowEnd,
	)

	if alreadyPlaced > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Area is not available")
	}

	_, err = o.db.Exec(
		`
			INSERT INTO yard_plans(
				block_id,
				slot_start, slot_end,
				row_start, row_end,
				size, height, type,
				slot_priority, row_priority, tier_priority
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		data.BlockId,
		data.SlotStart, data.SlotEnd,
		data.RowStart, data.RowEnd,
		data.Size, data.Height, data.Type,
		data.SlotPriority, data.RowPriority, data.TierPriority,
	)
	if err != nil {
		return err
	}

	return nil
}

func (o *YardPlanService) Delete(id string) error {
	_, err := o.db.Exec("DELETE FROM yard_plans WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (o *YardPlanService) Suggest(data *validation.SuggestionPayload) (*response.Suggestion, error) {
	plans := []response.YardPlan{}
	err := o.db.Select(
		&plans,
		`
			SELECT yard_plans.* FROM yard_plans
			JOIN blocks ON blocks.id = yard_plans.block_id
			WHERE yard_id = $1 AND type = $2
			ORDER BY yard_plans.id
		`,
		data.YardId,
		data.ContainerType,
	)

	if err != nil {
		return nil, err
	}

	if len(plans) == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Yard plan not found")
	}
	
	slotSize := 1
	if data.ContainerSize == 40 {
		slotSize = 2
	}
	
	type Slot struct {
		Slot int `db:"slot"`
		Row  int `db:"row"`
		Tier int `db:"tier"`
	}
	
	var suggestion *response.Suggestion
	
	for _, plan := range plans {
		if plan.SlotEnd - plan.SlotStart + 1 < slotSize {
			continue
		}
		
		block, err := o.blockService.GetById(plan.BlockId)
	
		if block == nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Block not found")
		}
		
		slots := []Slot{}
		err = o.db.Select(
			&slots,
			`
				SELECT slot, row, tier FROM placements
				WHERE block_id = $1
					AND slot >= $2 AND slot <= $3
					AND row >= $4 AND row <= $5
			`,
			plan.BlockId,
			plan.SlotStart,
			plan.SlotEnd,
			plan.RowStart,
			plan.RowEnd,
		)

		if err != nil {
			return nil, err
		}
		
		// Very inefficient, no time sorry
		outer:
		for slot := plan.SlotStart; slot <= plan.SlotEnd; slot++ {
			for row := plan.RowStart; row <= plan.RowEnd; row++ {
				for tier := 1; tier <= block.Tiers; tier++ {
					foundSlot := utils.ArrayFind(slots, func(t Slot) bool {
						if slotSize == 2 {
							return (t.Slot == slot || t.Slot == slot + 1) && t.Row == row && t.Tier == tier
						} else {
							return t.Slot == slot && t.Row == row && t.Tier == tier
						}
					})
					
					if foundSlot == nil {
						suggestion = &response.Suggestion{
							BlockId: block.Id,
							Slot: slot,
							Row: row,
							Tier: tier,
						}
						
						break outer
					}
				}
			}
		}
	}

	return suggestion, nil
}

func (o *YardPlanService) CheckAvailability(blockId int, slot int, slotSize int, row int, tier int) (bool, error) {
	var placedCount int
	var err error

	if slotSize == 1 {
		err = o.db.Get(
			&placedCount,
			`
				SELECT COUNT(*) FROM placements
				WHERE block_id = $1
					AND slot = $2
					AND row = $3
					AND tier = $4
			`,
			blockId,
			slot,
			row,
			tier,
		)
	} else if slotSize == 2 {
		err = o.db.Get(
			&placedCount,
			`
				SELECT COUNT(*) FROM placements
				WHERE block_id = $1
					AND slot = $2 OR slot = $3
					AND row = $4
					AND tier = $5
			`,
			blockId,
			slot,
			slot+1,
			row,
			tier,
		)
	}

	if err != nil {
		return false, err
	}

	return placedCount == 0, nil
}

func (o *YardPlanService) Place(data *validation.PlacementPayload) error {
	slotSize := 1
	if data.ContainerSize == 40 {
		slotSize = 2
	}

	var availablePlans int

	if slotSize == 1 {
		err := o.db.Get(
			&availablePlans,
			`
				SELECT COUNT(*) FROM yard_plans
				WHERE slot_start <= $1 AND slot_end >= $1 AND row_start <= $2 AND row_end >= $2
			`,
			data.Slot,
			data.Row,
		)
		if err != nil {
			return err
		}
	} else if slotSize == 2 {
		err := o.db.Get(
			&availablePlans,
			`
				SELECT COUNT(*) FROM yard_plans
				WHERE slot_start <= $1 AND slot_end >= $1 + 1 AND row_start <= $2 AND row_end >= $2
			`,
			data.Slot,
			data.Row,
		)
		if err != nil {
			return err
		}
	}

	if availablePlans == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No plan is available for the area")
	}

	available, err := o.CheckAvailability(data.BlockId, data.Slot, slotSize, data.Row, data.Tier)
	if err != nil {
		return err
	}

	if !available {
		return fiber.NewError(fiber.StatusBadRequest, "Slot not available")
	}

	tx, err := o.db.Beginx()

	_, err = tx.Exec(
		`
			INSERT INTO placements(
				container_id,
				container_size,
				container_height,
				container_type,
				block_id,
				slot,
				row,
				tier,
				is_head
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		data.ContainerId,
		data.ContainerSize,
		data.ContainerHeight,
		data.ContainerType,
		data.BlockId,
		data.Slot,
		data.Row,
		data.Tier,
		true,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if slotSize == 2 {
		_, err = tx.Exec(
			`
				INSERT INTO placements(
					container_id,
					container_size,
					container_height,
					container_type,
					block_id,
					slot,
					row,
					tier,
					is_head
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`,
			data.ContainerId,
			data.ContainerSize,
			data.ContainerHeight,
			data.ContainerType,
			data.BlockId,
			data.Slot+1,
			data.Row,
			data.Tier,
			false,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (o *YardPlanService) Pickup(data *validation.PickupPayload) error {
	res, err := o.db.Exec(
		`
			DELETE FROM placements WHERE container_id = $1
		`,
		data.ContainerId,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()

	if affected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No container found")
	}

	return nil
}
