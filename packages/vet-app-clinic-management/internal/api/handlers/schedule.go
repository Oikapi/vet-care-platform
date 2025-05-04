package handlers

import (
    "vet-app-clinic-management/internal/api/dto"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/service"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "log"
)

// ScheduleHandler обрабатывает HTTP-запросы для расписаний
type ScheduleHandler struct {
    svc *service.ScheduleService
}

// NewScheduleHandler создает новый экземпляр ScheduleHandler
func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{svc: svc}
}

// GetByDoctorID godoc
// @Summary Get schedules by doctor ID
// @Description Get all schedules for specific doctor
// @Tags schedules
// @Accept json
// @Produce json
// @Param doctorID path int true "Doctor ID"
// @Success 200 {array} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules/doctor/:doctorID [get]

func (h *ScheduleHandler) GetAll(c *gin.Context) {
    log.Println("Received request to get all schedules")
    schedules, err := h.svc.GetAll()
    if err != nil {
        log.Println("Failed to get all schedules:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Retrieved %d schedules", len(schedules))
    c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) GetByDoctorID(c *gin.Context) {
    log.Println("Received request to get schedules by doctor ID")
    doctorID := c.Param("doctorID")
    id, err := strconv.Atoi(doctorID)
    if err != nil {
        log.Println("Invalid doctor ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
        return
    }
    schedules, err := h.svc.GetByDoctorID(id)
    if err != nil {
        log.Println("Failed to get schedules by doctor ID:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Retrieved %d schedules for doctor ID %d", len(schedules), id)
    c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) Create(c *gin.Context) {
    var req dto.ScheduleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    if _, err := h.svc.GetDoctorByID(req.DoctorID); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
        return
    }

    schedule := &models.Schedule{
        DoctorID:  req.DoctorID,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
    }
    if err := h.svc.Create(schedule); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
    scheduleID, err := strconv.Atoi(c.Param("scheduleID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
        return
    }
    var req dto.ScheduleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    schedule := &models.Schedule{
        ID:        scheduleID,
        DoctorID:  req.DoctorID,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
    }
    if err := h.svc.Update(schedule); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
    scheduleID, err := strconv.Atoi(c.Param("scheduleID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
        return
    }
    if err := h.svc.Delete(scheduleID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}