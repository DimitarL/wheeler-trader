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

type RentHandler struct {
	st *storage.AppStorage
}

func NewRentHandler(st *storage.AppStorage) *RentHandler {
	return &RentHandler{st: st}
}

func (rnt *RentHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var rent entities.Rent
	err = json.Unmarshal(bodyBytes, &rent)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	err = rent.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	createdRent, err := rnt.st.CreateRent(rent)
	if err != nil {
		if common.IsForeignKeyError(err) {
			common.RespondWithErr(w, http.StatusBadRequest, errors.New("no vehicle with the specified id found"))
			return
		}
		if common.IsDuplicateKeyError(err) {
			common.RespondWithErr(w, http.StatusBadRequest, errors.New("rent for vehicle with the specified id already exists"))
			return
		}

		err := fmt.Errorf("error creating rent %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, createdRent, http.StatusCreated)
}

func (rnt *RentHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	params, err := extractRentFilterParameters(r.URL.Query())
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	rents, err := rnt.st.ListAllRents(params)
	if err != nil {
		err := fmt.Errorf("error getting all rents %w", err)
		common.RespondWithInternalErr(w, err)
		return
	}

	common.RespondWithJson(w, rents, http.StatusOK)
}

func extractRentFilterParameters(query url.Values) (map[string]interface{}, error) {
	return common.ExtractQueryParameters(query, []string{"type", "make", "model"}, []string{"minHP", "maxHP", "minDP", "maxDP", "minWP", "maxWP", "minMP", "maxMP"})
}

func (rnt *RentHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	rent, err := rnt.st.GetRentByVehicleId(currId)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}
	if rent == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("rental of vehicle with id %d not found", currId))
		return
	}

	common.RespondWithJson(w, rent, http.StatusOK)
}

func (rnt *RentHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	var rent entities.Rent
	err = json.Unmarshal(bodyBytes, &rent)
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}
	err = rent.Validate()
	if err != nil {
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	params := mux.Vars(r)
	rent.VehicleID, err = strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	updatedRent, err := rnt.st.UpdateRentByVehicleId(rent)
	if err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}
	if updatedRent == nil {
		common.RespondWithErr(w, http.StatusNotFound, fmt.Errorf("rent vehicle with id %d not found", rent.VehicleID))
		return
	}

	common.RespondWithJson(w, updatedRent, http.StatusOK)
}

func (rnt *RentHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	currId, err := strconv.Atoi(params["id"])
	if err != nil {
		err = fmt.Errorf("parameter 'id': %w", err)
		common.RespondWithErr(w, http.StatusBadRequest, err)
		return
	}

	if err = rnt.st.DeleteRentByVehicleId(currId); err != nil {
		common.RespondWithInternalErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
