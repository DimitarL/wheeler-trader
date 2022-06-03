package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DimitarL/wheeler-trader/server/common"
	"github.com/DimitarL/wheeler-trader/server/pkg/entities"
	"github.com/DimitarL/wheeler-trader/server/pkg/storage"
	"github.com/gorilla/mux"
)

type SaleHandler struct {
	st *storage.AppStorage
}

func NewSaleHandler(st *storage.AppStorage) *SaleHandler {
	return &SaleHandler{st: st}
}

func (s *SaleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var sale entities.Sale
	err = json.Unmarshal(bodyBytes, &sale)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	err = sale.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	createdSale, err := s.st.CreateSale(sale)
	if err != nil {
		if common.IsForeignKeyError(err) {
			common.RespondWithErr(w, http.StatusBadRequest, errors.New("no vehicle with the specified id found"))
			return
		}
		if common.IsDuplicateKeyError(err) {
			common.RespondWithErr(w, http.StatusBadRequest, errors.New("sale for vehicle with the specified id already exists"))
			return
		}

		err := fmt.Errorf("error creating sale %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, createdSale, http.StatusCreated)
}

func (s *SaleHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	params, err := extractSaleFilterParameters(r.URL.Query())
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	sales, err := s.st.ListAllSales(params)
	if err != nil {
		err := fmt.Errorf("error getting all sales %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, sales, http.StatusOK)
}

func extractSaleFilterParameters(query url.Values) (map[string]interface{}, error) {
	return common.ExtractQueryParameters(query, []string{"type", "make", "model"}, []string{"minHP", "maxHP", "minPrice", "maxPrice"})
}

func (s *SaleHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	sale, err := s.st.GetSaleByVehicleId(currId)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}
	if sale == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("sale of vehicle with id %d not found", currId))
		return
	}

	common.RespondWithJson(w, sale, http.StatusOK)
}

func (s *SaleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var sale entities.Sale
	err = json.Unmarshal(bodyBytes, &sale)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}
	err = sale.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	params := mux.Vars(r)
	sale.VehicleID, err = strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	updatedSale, err := s.st.UpdateSaleByVehicleId(sale)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}
	if updatedSale == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("sale vehicle with id %d not found", sale.VehicleID))
		return
	}

	common.RespondWithJson(w, updatedSale, http.StatusOK)
}

func (s *SaleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	if err = s.st.DeleteSaleByVehicleId(currId); err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
