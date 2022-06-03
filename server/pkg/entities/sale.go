package entities

import (
	"errors"
)

type Sale struct {
	VehicleID int `json:"vehicle_id"`
	Price     int `json:"price"`
}

func (s Sale) Validate() error {
	if s.Price <= 0 {
		return errors.New("property 'price' should be positive number")
	}

	return nil
}
