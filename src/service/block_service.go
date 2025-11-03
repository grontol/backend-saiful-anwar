package service

import (
	"yard_plan/src/response"
	"yard_plan/src/validation"

	"github.com/jmoiron/sqlx"
)

type BlockService struct {
	db *sqlx.DB
}

func NewBlockService(
	db *sqlx.DB,
) *BlockService {
	return &BlockService{
		db: db,
	}
}

func (o *BlockService) List() ([]response.Block, error) {
	data := []response.Block{}
	err := o.db.Select(&data, "SELECT * FROM blocks ORDER BY id")
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (o *BlockService) ListByYard(yardId string) ([]response.Block, error) {
	data := []response.Block{}
	err := o.db.Select(&data, "SELECT * FROM blocks WHERE yard_id = $1 ORDER BY id", yardId)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (o *BlockService) GetById(id int) (*response.Block, error) {
	data := new(response.Block)
	err := o.db.Get(data, "SELECT * FROM blocks WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (o *BlockService) Create(data *validation.BlockPayload) error {
	_, err := o.db.Exec(
		"INSERT INTO blocks (yard_id, name, slots, rows, tiers) VALUES ($1, $2, $3, $4, $5)",
		data.YardId,
		data.Name,
		data.Slots,
		data.Rows,
		data.Tiers,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (o *BlockService) Edit(id string, data *validation.BlockPayload) error {
	_, err := o.db.Exec(
		"UPDATE blocks SET yard_id = $1, name = $2, slots = $3, rows = $4, tiers = $5 WHERE id = $6",
		data.YardId,
		data.Name,
		data.Slots,
		data.Rows,
		data.Tiers,
		id,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (o *BlockService) Delete(id string) error {
	_, err := o.db.Exec("DELETE FROM blocks WHERE id = $1", id)
	if err != nil {
		return err
	}
	
	return nil
}