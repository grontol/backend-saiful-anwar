package service

import (
	"yard_plan/src/response"

	"github.com/jmoiron/sqlx"
)

type PlacementService struct {
	db *sqlx.DB
}

func NewPlacementService(
	db *sqlx.DB,
) *PlacementService {
	return &PlacementService{
		db: db,
	}
}

func (o *PlacementService) List() ([]response.Placement, error) {
	data := []response.Placement{}
	err := o.db.Select(&data, "SELECT * FROM placements ORDER BY id")
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (o *PlacementService) ListByBlock(blockId string) ([]response.YardPlan, error) {
	data := []response.YardPlan{}
	err := o.db.Select(&data, "SELECT * FROM placements WHERE block_id = $1 ORDER BY id", blockId)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}