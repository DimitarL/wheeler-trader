package entities

import (
	"errors"
)

type Rent struct {
	VehicleID    int `json:"vehicle_id"`
	DailyPrice   int `json:"daily_price"`
	WeeklyPrice  int `json:"weekly_price"`
	MonthlyPrice int `json:"monthly_price"`
}

func (r Rent) Validate() error {
	if r.DailyPrice <= 0 {
		return errors.New("property 'dailyPrice' should be positive number")
	}
	if r.WeeklyPrice <= 0 {
		return errors.New("property 'weeklyPrice' should be positive number")
	}
	if r.MonthlyPrice <= 0 {
		return errors.New("property 'monthlyPrice' should be positive number")
	}

	return nil
}
