package handler

import (
	"encoding/json"
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

type VehicleHandler struct {
	st *storage.AppStorage
}

func NewVehicleHandler(st *storage.AppStorage) *VehicleHandler {
	return &VehicleHandler{st: st}
}

func (h *VehicleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var vehicle entities.Vehicle
	err = json.Unmarshal(bodyBytes, &vehicle)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	err = vehicle.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	createdVehicle, err := h.st.CreateVehicle(vehicle)
	if err != nil {
		err := fmt.Errorf("error creating vehicle %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, createdVehicle, http.StatusCreated)
}

func (h *VehicleHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	params, err := extractFilterParameters(r.URL.Query())
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	vehicles, err := h.st.ListVehicles(params)
	if err != nil {
		err := fmt.Errorf("error getting all vehicles %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, vehicles, http.StatusOK)
}

func extractFilterParameters(query url.Values) (map[string]interface{}, error) {
	return common.ExtractQueryParameters(query, []string{"type", "make", "model"}, []string{"minHP", "maxHP"})
}

func (h *VehicleHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	vehicle, err := h.st.GetVehicleById(currId)
	if err != nil {
		err = fmt.Errorf("error getting vehicle with id %d: %w", currId, err)
		common.RespondWithInternalErr(w, err)
		return
	}
	if vehicle == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("vehicle with id %d not found", currId))
		return
	}

	common.RespondWithJson(w, vehicle, http.StatusOK)
}

func (h *VehicleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var vehicle entities.Vehicle
	err = json.Unmarshal(bodyBytes, &vehicle)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}
	err = vehicle.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	params := mux.Vars(r)
	vehicle.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	updatedVehicle, err := h.st.UpdateVehicleById(vehicle)
	if err != nil {
		err = fmt.Errorf("error updating vehicle with id %d: %w", vehicle.ID, err)
		common.RespondWithInternalErr(w, err)
		return
	}
	if updatedVehicle == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("vehicle with id %d not found", vehicle.ID))
		return
	}

	common.RespondWithJson(w, updatedVehicle, http.StatusOK)
}

func (h *VehicleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	if err = h.st.DeleteVehicleById(currId); err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
