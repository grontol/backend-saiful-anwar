package service

import (
	"yard_plan/src/response"
	"yard_plan/src/validation"

	"github.com/jmoiron/sqlx"
)

type YardService struct {
	db *sqlx.DB
}

func NewYardService(
	db *sqlx.DB,
) *YardService {
	return &YardService{
		db: db,
	}
}

func (o *YardService) List() ([]response.Yard, error) {
	data := []response.Yard{}
	err := o.db.Select(&data, "SELECT * FROM yards ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (o *YardService) Create(data *validation.YardPayload) error {
	_, err := o.db.Exec("INSERT INTO yards (name, description) VALUES ($1, $2)", data.Name, data.Description)
	if err != nil {
		return err
	}
	
	return nil
}

func (o *YardService) Edit(id string, data *validation.YardPayload) error {
	_, err := o.db.Exec("UPDATE yards SET name = $1, description = $2 WHERE id = $3", data.Name, data.Description, id)
	if err != nil {
		return err
	}
	
	return nil
}

func (o *YardService) Delete(id string) error {
	_, err := o.db.Exec("DELETE FROM yards WHERE id = $1", id)
	if err != nil {
		return err
	}
	
	return nil
}