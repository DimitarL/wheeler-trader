package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimitarL/wheeler-trader/server/pkg/entities"
	pgx "github.com/jackc/pgx/v4"
)

type InventoryResponse struct {
	Vehicle      entities.Vehicle `json:"vehicle"`
	Price        *int             `json:"price,omitempty"`
	DailyPrice   *int             `json:"daily_price,omitempty"`
	WeeklyPrice  *int             `json:"weekly_price,omitempty"`
	MonthlyPrice *int             `json:"monthly_price,omitempty"`
}

const (
	selectPartOfQuery = `SELECT id, type, make, model, horsepower, price, daily_price, weekly_price, monthly_price
											FROM vehicle LEFT JOIN sale
											ON vehicle.id = sale.vehicle_id
											LEFT JOIN rent
											ON vehicle.id = rent.vehicle_id`
	selectVehiclesQuery = `SELECT id, type, make, model, horsepower
											FROM vehicle`
	getByIdInventoryQuery = `SELECT id, type, make, model, horsepower, price, daily_price, weekly_price, monthly_price
											FROM vehicle LEFT JOIN sale
											ON vehicle.id = sale.vehicle_id
											LEFT JOIN rent
											ON vehicle.id = rent.vehicle_id
											WHERE id = $1;`
)

func (a AppStorage) ListInventory(params map[string]interface{}) ([]interface{}, error) {
	vehicles := make([]interface{}, 0)

	getAllQuery, sqlParams := BuildGetAllQuery(selectPartOfQuery, params)

	rows, err := a.conn.Query(context.Background(), getAllQuery, sqlParams...)
	if err != nil {
		return vehicles, err
	}
	defer rows.Close()

	for rows.Next() {
		var data InventoryResponse

		err := rows.Scan(expandInventory(&data)...)
		if err != nil {
			return vehicles, fmt.Errorf("failed to scan row: %w", err)
		}

		vehicles = append(vehicles, data)
	}

	if err := rows.Err(); err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

func (a AppStorage) ListUnassignedVehicles(params map[string]interface{}) ([]entities.Vehicle, error) {
	vehicles := []entities.Vehicle{}
	var getUnassignedQuery, andQueryPart string
	var sqlParams []interface{}

	if len(params) != 0 {
		getUnassignedQuery, sqlParams = BuildGetAllQuery(selectVehiclesQuery, params)
		andQueryPart = `AND id NOT IN (
			SELECT vehicle_id FROM sale
			UNION
			SELECT vehicle_id FROM rent
		)`
	} else {
		getUnassignedQuery = selectVehiclesQuery
		andQueryPart = ` WHERE id NOT IN (
			SELECT vehicle_id FROM sale
			UNION
			SELECT vehicle_id FROM rent
		)`
	}

	rows, err := a.conn.Query(context.Background(), getUnassignedQuery+andQueryPart, sqlParams...)
	if err != nil {
		return vehicles, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle entities.Vehicle

		err := rows.Scan(expandVehicle(&vehicle)...)
		if err != nil {
			return vehicles, fmt.Errorf("failed to scan row: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

func (a AppStorage) GetInventoryById(id int) (*InventoryResponse, error) {
	var vehicle InventoryResponse

	err := a.conn.QueryRow(context.Background(), getByIdInventoryQuery, id).Scan(expandInventory(&vehicle)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &vehicle, err
}

func expandInventory(veh *InventoryResponse) []interface{} {
	return []interface{}{&veh.Vehicle.ID, &veh.Vehicle.Type, &veh.Vehicle.Make, &veh.Vehicle.Model, &veh.Vehicle.HorsePower,
		&veh.Price, &veh.DailyPrice, &veh.WeeklyPrice, &veh.MonthlyPrice}
}
