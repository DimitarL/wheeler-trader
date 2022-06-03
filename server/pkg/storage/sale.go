package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimitarL/wheeler-trader/server/pkg/entities"
	pgx "github.com/jackc/pgx/v4"
)

const (
	insertSaleQuery = `INSERT INTO sale(vehicle_id, price) VALUES ($1, $2)
										RETURNING vehicle_id, price;`
	getByIdSaleQuery = `SELECT v.id, v.type, v.make, v.model, v.horsepower, s.price
											FROM sale s JOIN vehicle v
											ON v.id = s.vehicle_id
											WHERE s.vehicle_id = $1;`
	updateSaleQuery = `UPDATE sale SET price = $2 WHERE vehicle_id = $1
										RETURNING vehicle_id, price;`
	deleteSaleQuery = `DELETE FROM sale WHERE vehicle_id = $1;`
)

type SaleResponse struct {
	Vehicle entities.Vehicle `json:"vehicle"`
	Price   int              `json:"price"`
}

func (a AppStorage) CreateSale(sale entities.Sale) (*entities.Sale, error) {
	var createdSale entities.Sale

	return &createdSale, a.conn.QueryRow(context.Background(), insertSaleQuery, sale.VehicleID, sale.Price).Scan(expandSaleVehicle(&createdSale)...)
}

func (a AppStorage) ListAllSales(params map[string]interface{}) ([]SaleResponse, error) {
	sales := []SaleResponse{}

	selectPartOfQuery := `SELECT id, type, make, model, horsepower, price
												FROM vehicle JOIN sale
												ON vehicle.id = sale.vehicle_id`
	getAllQuery, sqlParams := BuildGetAllQuery(selectPartOfQuery, params)

	rows, err := a.conn.Query(context.Background(), getAllQuery, sqlParams...)
	if err != nil {
		return sales, err
	}
	defer rows.Close()

	for rows.Next() {
		var data SaleResponse

		err := rows.Scan(expandCompleteSaleEntry(&data)...)
		if err != nil {
			return sales, fmt.Errorf("failed to scan row: %w", err)
		}

		sales = append(sales, data)
	}

	if err := rows.Err(); err != nil {
		return sales, err
	}

	return sales, nil
}

func (a AppStorage) GetSaleByVehicleId(id int) (*SaleResponse, error) {
	var responseData SaleResponse

	err := a.conn.QueryRow(context.Background(), getByIdSaleQuery, id).Scan(expandCompleteSaleEntry(&responseData)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &responseData, nil
}

func (a AppStorage) UpdateSaleByVehicleId(updatedVehicle entities.Sale) (*entities.Sale, error) {
	var sale entities.Sale

	err := a.conn.QueryRow(context.Background(), updateSaleQuery, expandSaleVehicle(&updatedVehicle)...).
		Scan(expandSaleVehicle(&sale)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return &sale, nil
}

func (a AppStorage) DeleteSaleByVehicleId(id int) error {
	_, err := a.conn.Exec(context.Background(), deleteSaleQuery, id)

	return err
}

func expandSaleVehicle(veh *entities.Sale) []interface{} {
	return []interface{}{&veh.VehicleID, &veh.Price}
}

func expandCompleteSaleEntry(data *SaleResponse) []interface{} {
	return []interface{}{&data.Vehicle.ID, &data.Vehicle.Type, &data.Vehicle.Make, &data.Vehicle.Model,
		&data.Vehicle.HorsePower, &data.Price}
}
