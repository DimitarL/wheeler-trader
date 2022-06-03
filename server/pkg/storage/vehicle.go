package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DimitarL/wheeler-trader/server/pkg/entities"
	"github.com/jackc/pgx/v4"
)

const (
	insertQuery = `INSERT INTO vehicle(type, make, model, horsepower) VALUES ($1, $2, $3, $4)
								RETURNING vehicle.id, vehicle.type, vehicle.make, vehicle.model, vehicle.horsepower;`
	getByIdQuery = "SELECT id, type, make, model, horsepower FROM vehicle WHERE id = $1;"
	updateQuery  = `UPDATE vehicle SET type = $2, make = $3, model = $4, horsepower = $5 WHERE id = $1
									RETURNING vehicle.id, vehicle.type, vehicle.make, vehicle.model, vehicle.horsepower;`
	deleteQuery = "DELETE FROM vehicle WHERE id = $1;"
)

func (a AppStorage) CreateVehicle(veh entities.Vehicle) (*entities.Vehicle, error) {
	var vehInDB entities.Vehicle

	return &vehInDB, a.conn.QueryRow(context.Background(), insertQuery, veh.Type, veh.Make, veh.Model, veh.HorsePower).Scan(expandVehicle(&vehInDB)...)
}

func (a AppStorage) ListVehicles(params map[string]interface{}) ([]entities.Vehicle, error) {
	vehicles := []entities.Vehicle{}
	selectPartOfQuery := "SELECT id, type, make, model, horsepower FROM vehicle"
	getAllQuery, sqlParams := BuildGetAllQuery(selectPartOfQuery, params)

	rows, err := a.conn.Query(context.Background(), getAllQuery, sqlParams...)
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

func (a AppStorage) GetVehicleById(id int) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle

	err := a.conn.QueryRow(context.Background(), getByIdQuery, id).Scan(expandVehicle(&vehicle)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &vehicle, err
}

func (a AppStorage) UpdateVehicleById(updatedVehicle entities.Vehicle) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle

	err := a.conn.QueryRow(context.Background(), updateQuery, expandVehicle(&updatedVehicle)...).
		Scan(expandVehicle(&vehicle)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &vehicle, err
}

func (a AppStorage) DeleteVehicleById(id int) error {
	_, err := a.conn.Exec(context.Background(), deleteQuery, id)

	return err
}

func BuildGetAllQuery(selectPartOfQuery string, params map[string]interface{}) (string, []interface{}) {
	wherePartOfQuery := ""
	var sqlParams []interface{}

	if len(params) != 0 {
		wherePartOfQuery, sqlParams = buildWhereClause(params)
	}
	return selectPartOfQuery + wherePartOfQuery, sqlParams
}

func buildWhereClause(params map[string]interface{}) (string, []interface{}) {
	return newWhereClauseBuilder(getColumnInfo).Build(params)
}

func expandVehicle(veh *entities.Vehicle) []interface{} {
	return []interface{}{&veh.ID, &veh.Type, &veh.Make, &veh.Model, &veh.HorsePower}
}

func getColumnInfo(param string) columnInfo {
	info := columnInfo{}
	param = strings.ToLower(param)

	if strings.HasPrefix(param, "min") || strings.Contains(param, "max") {
		info.dbType = numeric
		param = param[3:]
	} else {
		info.dbType = character
	}

	switch param {
	case "type":
		info.name = "type"
	case "make":
		info.name = "make"
	case "model":
		info.name = "model"
	case "hp":
		info.name = "horsepower"
	case "price":
		info.name = "price"
	case "dp":
		info.name = "daily_price"
	case "wp":
		info.name = "weekly_price"
	case "mp":
		info.name = "monthly_price"
	default:
		panic(fmt.Sprintf("could not handle param \"%s\"", param))
	}

	return info
}
