package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimitarL/wheeler-trader/server/pkg/entities"
	pgx "github.com/jackc/pgx/v4"
)

const (
	insertRentQuery = `INSERT INTO rent(vehicle_id, daily_price, weekly_price, monthly_price)
										VALUES ($1, $2, $3, $4)
										RETURNING vehicle_id, daily_price, weekly_price, monthly_price;`
	getByIdRentQuery = `SELECT v.id, v.type, v.make, v.model, v.horsepower, r.daily_price, r.weekly_price, r.monthly_price
											FROM rent r JOIN vehicle v
											ON v.id = r.vehicle_id
											WHERE r.vehicle_id = $1;`
	updateRentQuery = `UPDATE rent SET daily_price = $2, weekly_price = $3, monthly_price = $4
										WHERE vehicle_id = $1
										RETURNING vehicle_id, daily_price, weekly_price, monthly_price;`
	deleteRentQuery = `DELETE FROM rent WHERE vehicle_id = $1;`
)

type RentResponse struct {
	Vehicle      entities.Vehicle `json:"vehicle"`
	DailyPrice   int              `json:"daily_price"`
	WeeklyPrice  int              `json:"weekly_price"`
	MonthlyPrice int              `json:"monthly_price"`
}

func (a AppStorage) CreateRent(rent entities.Rent) (*entities.Rent, error) {
	var createdRent entities.Rent

	return &createdRent, a.conn.QueryRow(context.Background(), insertRentQuery, rent.VehicleID, rent.DailyPrice,
		rent.WeeklyPrice, rent.MonthlyPrice).Scan(expandRentVehicle(&createdRent)...)
}

func (a AppStorage) ListAllRents(params map[string]interface{}) ([]RentResponse, error) {
	rents := []RentResponse{}

	selectPartOfQuery := `SELECT id, type, make, model, horsepower, daily_price, weekly_price, monthly_price
												FROM vehicle JOIN rent
												ON vehicle.id = rent.vehicle_id`
	getAllQuery, sqlParams := BuildGetAllQuery(selectPartOfQuery, params)

	rows, err := a.conn.Query(context.Background(), getAllQuery, sqlParams...)
	if err != nil {
		return rents, err
	}
	defer rows.Close()

	for rows.Next() {
		var data RentResponse

		err := rows.Scan(expandCompleteRentEntry(&data)...)
		if err != nil {
			return rents, fmt.Errorf("failed to scan row: %w", err)
		}

		rents = append(rents, data)
	}

	if err := rows.Err(); err != nil {
		return rents, err
	}

	return rents, nil
}

func (a AppStorage) GetRentByVehicleId(id int) (*RentResponse, error) {
	var responseData RentResponse

	err := a.conn.QueryRow(context.Background(), getByIdRentQuery, id).Scan(expandCompleteRentEntry(&responseData)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &responseData, nil
}

func (a AppStorage) UpdateRentByVehicleId(updatedVehicle entities.Rent) (*entities.Rent, error) {
	var rent entities.Rent

	err := a.conn.QueryRow(context.Background(), updateRentQuery, expandRentVehicle(&updatedVehicle)...).
		Scan(expandRentVehicle(&rent)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &rent, nil
}

func (a AppStorage) DeleteRentByVehicleId(id int) error {
	_, err := a.conn.Exec(context.Background(), deleteRentQuery, id)

	return err
}

func expandRentVehicle(rent *entities.Rent) []interface{} {
	return []interface{}{&rent.VehicleID, &rent.DailyPrice, &rent.WeeklyPrice, &rent.MonthlyPrice}
}

func expandCompleteRentEntry(rent *RentResponse) []interface{} {
	return []interface{}{&rent.Vehicle.ID, &rent.Vehicle.Type, &rent.Vehicle.Make, &rent.Vehicle.Model,
		&rent.Vehicle.HorsePower, &rent.DailyPrice, &rent.WeeklyPrice, &rent.MonthlyPrice}
}
