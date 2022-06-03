package entities

import (
	"errors"
)

type Vehicle struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Make       string `json:"make"`
	Model      string `json:"model"`
	HorsePower int    `json:"horsePower"`
}

func (v Vehicle) Validate() error {
	if v.Type != "car" && v.Type != "motorbike" && v.Type != "truck" {
		return errors.New("property 'type' should be either car, motorbike or truck")
	}
	if v.Make == "" {
		return errors.New("property 'make' should not be empty")
	}
	if v.Model == "" {
		return errors.New("property 'model' should not be empty")
	}
	if v.HorsePower <= 0 {
		return errors.New("property 'horsePower' should be positive number")
	}

	return nil
}
