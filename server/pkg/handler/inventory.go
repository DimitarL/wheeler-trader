package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DimitarL/wheeler-trader/server/common"
	"github.com/DimitarL/wheeler-trader/server/pkg/storage"
	"github.com/gorilla/mux"
)

type InventoryHandler struct {
	st *storage.AppStorage
}

func NewInventoryHandler(st *storage.AppStorage) *InventoryHandler {
	return &InventoryHandler{st: st}
}

func (h *InventoryHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	params, err := extractAssignedFilterParameters(r.URL.Query())
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	vehicles, err := h.st.ListInventory(params)
	if err != nil {
		err := fmt.Errorf("error getting all vehicles %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, vehicles, http.StatusOK)
}

func extractAssignedFilterParameters(query url.Values) (map[string]interface{}, error) {
	return common.ExtractQueryParameters(query, []string{"type", "make", "model"}, []string{"minHP", "maxHP", "minPrice", "maxPrice", "minDP", "maxDP", "minWP", "maxWP", "minMP", "maxMP"})
}

func (h *InventoryHandler) GetUnassignedHandler(w http.ResponseWriter, r *http.Request) {
	params, err := extractUnassignedFilterParameters(r.URL.Query())
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	vehicles, err := h.st.ListUnassignedVehicles(params)
	if err != nil {
		err := fmt.Errorf("error getting all unassigned vehicles %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, vehicles, http.StatusOK)
}

func extractUnassignedFilterParameters(query url.Values) (map[string]interface{}, error) {
	return common.ExtractQueryParameters(query, []string{"type", "make", "model"}, []string{"minHP", "maxHP"})
}

func (h *InventoryHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	vehicle, err := h.st.GetInventoryById(currId)
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
