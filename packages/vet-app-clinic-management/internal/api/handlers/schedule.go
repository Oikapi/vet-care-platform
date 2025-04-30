package handlers

import (
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

// ServeHTTP реализует интерфейс http.Handler
func (h *ScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByDoctor(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByDoctor godoc
// @Summary Get schedules by doctor ID
// @Description Get all schedules for specific doctor
// @Tags schedules
// @Accept json
// @Produce json
// @Param doctorID path int true "Doctor ID"
// @Success 200 {array} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/schedules/{doctorID} [get]
func (h *ScheduleHandler) GetByDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorIDStr := vars["doctorID"]
	
	doctorID, err := strconv.ParseUint(doctorIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	schedules, err := h.svc.GetByDoctor(uint(doctorID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, schedules)
}

