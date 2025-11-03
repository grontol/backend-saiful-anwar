package validation

type YardPayload struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:""`
}

type BlockPayload struct {
	YardId int    `json:"yard_id" validate:"required"`
	Name   string `json:"name" validate:"required,max=255"`
	Slots  int    `json:"slots" validate:"required"`
	Rows   int    `json:"rows" validate:"required"`
	Tiers  int    `json:"tiers" validate:"required"`
}

type YardPlanPayload struct {
	BlockId      int     `json:"block_id" validate:"required"`
	SlotStart    int     `json:"slot_start" validate:"required"`
	SlotEnd      int     `json:"slot_end" validate:"required"`
	RowStart     int     `json:"row_start" validate:"required"`
	RowEnd       int     `json:"row_end" validate:"required"`
	Size         int     `json:"size" validate:"required,oneof=20 40"`
	Height       float32 `json:"height" validate:"required"`
	Type         string  `json:"type" validate:"required,oneof=DRY REEFER OPEN_TOP FLAT_RACK"`
	SlotPriority int     `json:"slot_priority" validate:""`
	RowPriority  int     `json:"row_priority" validate:""`
	TierPriority int     `json:"tier_priority" validate:""`
}

type PlacementPayload struct {
	ContainerId     string  `json:"container_id" validate:"required"`
	ContainerSize   int     `json:"container_size" validate:"required"`
	ContainerHeight float32 `json:"container_height" validate:"required"`
	ContainerType   string  `json:"container_type" validate:"required,oneof=DRY REEFER OPEN_TOP FLAT_RACK"`
	BlockId         int     `json:"block_id" validate:"required"`
	Slot            int     `json:"slot" validate:"required"`
	Row             int     `json:"row" validate:"required"`
	Tier            int     `json:"tier" validate:"required"`
}

type SuggestionPayload struct {
	YardId          int     `json:"yard_id" validate:"required"`
	ContainerId     string  `json:"container_id" validate:"required"`
	ContainerSize   int     `json:"container_size" validate:"required"`
	ContainerHeight float32 `json:"container_height" validate:"required"`
	ContainerType   string  `json:"container_type" validate:"required,oneof=DRY REEFER OPEN_TOP FLAT_RACK"`
}

type PickupPayload struct {
	ContainerId     string  `json:"container_id" validate:"required"`
}
