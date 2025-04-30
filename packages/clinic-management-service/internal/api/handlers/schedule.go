package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "clinic-management-service/internal/service"
    "github.com/gorilla/mux"
)

type ScheduleHandler struct {
    svc *service.ScheduleService
}

func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{svc: svc}
}

func (h *ScheduleHandler) GetByDoctor(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doctorIDStr := vars["doctorID"]
    // Конвертируем doctorID из строки в uint
    doctorID, err := strconv.ParseUint(doctorIDStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }
    // Используем doctorID
    schedules, err := h.svc.GetByDoctor(uint(doctorID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(schedules)
}