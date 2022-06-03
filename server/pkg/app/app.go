package app

import (
	"fmt"
	"net/http"

	"github.com/DimitarL/wheeler-trader/server/pkg/handler"
	"github.com/DimitarL/wheeler-trader/server/pkg/storage"
	"github.com/gorilla/mux"
)

type Application struct {
	Vehicles          *handler.VehicleHandler
	SaleVehicles      *handler.SaleHandler
	RentVehicles      *handler.RentHandler
	InventoryVehicles *handler.InventoryHandler
}

func NewApplication() *Application {
	st := storage.NewAppStorage()
	return &Application{
		Vehicles:          handler.NewVehicleHandler(st),
		SaleVehicles:      handler.NewSaleHandler(st),
		RentVehicles:      handler.NewRentHandler(st),
		InventoryVehicles: handler.NewInventoryHandler(st),
	}
}

func (a *Application) Start(host string, port int) error {
	router := a.createRouter()

	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router)
}

func (a *Application) createRouter() *mux.Router {
	router := mux.NewRouter()

	vehicleRouter := router.PathPrefix("/vehicle").Subrouter()
	vehicleRouter.HandleFunc("", a.Vehicles.CreateHandler).Methods("POST")
	vehicleRouter.HandleFunc("", a.Vehicles.GetAllHandler).Methods("GET")
	vehicleRouter.HandleFunc("/{id}", a.Vehicles.GetByIdHandler).Methods("GET")
	vehicleRouter.HandleFunc("/{id}", a.Vehicles.UpdateHandler).Methods("PUT")
	vehicleRouter.HandleFunc("/{id}", a.Vehicles.DeleteHandler).Methods("DELETE")

	saleRouter := router.PathPrefix("/sale").Subrouter()
	saleRouter.HandleFunc("", a.SaleVehicles.CreateHandler).Methods("POST")
	saleRouter.HandleFunc("", a.SaleVehicles.GetAllHandler).Methods("GET")
	saleRouter.HandleFunc("/{id}", a.SaleVehicles.GetByIdHandler).Methods("GET")
	saleRouter.HandleFunc("/{id}", a.SaleVehicles.UpdateHandler).Methods("PUT")
	saleRouter.HandleFunc("/{id}", a.SaleVehicles.DeleteHandler).Methods("DELETE")

	rentRouter := router.PathPrefix("/rent").Subrouter()
	rentRouter.HandleFunc("", a.RentVehicles.CreateHandler).Methods("POST")
	rentRouter.HandleFunc("", a.RentVehicles.GetAllHandler).Methods("GET")
	rentRouter.HandleFunc("/{id}", a.RentVehicles.GetByIdHandler).Methods("GET")
	rentRouter.HandleFunc("/{id}", a.RentVehicles.UpdateHandler).Methods("PUT")
	rentRouter.HandleFunc("/{id}", a.RentVehicles.DeleteHandler).Methods("DELETE")

	inventoryRouter := router.PathPrefix("/inventory").Subrouter()
	inventoryRouter.HandleFunc("/all", a.InventoryVehicles.GetAllHandler).Methods("GET")
	inventoryRouter.HandleFunc("/unassigned", a.InventoryVehicles.GetUnassignedHandler).Methods("GET")
	inventoryRouter.HandleFunc("/{id}", a.InventoryVehicles.GetByIdHandler).Methods("GET")

	return router
}
