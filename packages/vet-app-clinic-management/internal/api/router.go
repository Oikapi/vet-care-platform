package api

import (
    "clinic-management-service/internal/api/handlers"
    "github.com/gorilla/mux"
)

// SetupRouter настраивает маршруты для API.
func SetupRouter(scheduleHandler, inventoryHandler interface{}) *mux.Router {
    r := mux.NewRouter()
    // Приведение типов к нужным хендлерам
    sh := scheduleHandler.(*handlers.ScheduleHandler)
    ih := inventoryHandler.(*handlers.InventoryHandler)

    r.HandleFunc("/schedules/doctor/{doctorID}", sh.GetByDoctor).Methods("GET")
    r.HandleFunc("/inventory", ih.GetAll).Methods("GET")
    r.HandleFunc("/inventory/autoorder", ih.AutoOrder).Methods("POST")
    return r
}