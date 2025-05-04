package handlers

import (
    "vet-app-clinic-management/internal/api/dto"
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

// GetAll godoc
// @Summary Get all schedules
// @Description Get all schedules for a clinic
// @Tags schedules
// @Accept json
// @Produce json
// @Param clinicID path int true "Clinic ID"
// @Success 200 {array} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules [get]
func (h *ScheduleHandler) GetAll(c *gin.Context) {
    log.Println("Received request to get all schedules")
    clinicID := c.Param("clinicID")
    id, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    schedules, err := h.svc.GetAll(id)
    if err != nil {
        log.Println("Failed to get all schedules:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Retrieved %d schedules for clinic ID %d", len(schedules), id)
    c.JSON(http.StatusOK, schedules)
}

// GetByDoctorID godoc
// @Summary Get schedules by doctor ID
// @Description Get all schedules for specific doctor
// @Tags schedules
// @Accept json
// @Produce json
// @Param clinicID path int true "Clinic ID"
// @Param doctorID path int true "Doctor ID"
// @Success 200 {array} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules/doctor/:doctorID [get]
func (h *ScheduleHandler) GetByDoctorID(c *gin.Context) {
    log.Println("Received request to get schedules by doctor ID")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    doctorID := c.Param("doctorID")
    id, err := strconv.Atoi(doctorID)
    if err != nil {
        log.Println("Invalid doctor ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
        return
    }
    schedules, err := h.svc.GetByDoctorID(clinicIDInt, id)
    if err != nil {
        log.Println("Failed to get schedules by doctor ID:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Retrieved %d schedules for doctor ID %d and clinic ID %d", len(schedules), id, clinicIDInt)
    c.JSON(http.StatusOK, schedules)
}

// Create godoc
// @Summary Create a new schedule
// @Description Create a new schedule for a doctor
// @Tags schedules
// @Accept json
// @Produce json
// @Param clinicID path int true "Clinic ID"
// @Param schedule body dto.ScheduleRequest true "Schedule data"
// @Success 201 {object} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
    log.Println("Received request to create schedule")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    var req dto.ScheduleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    if _, err := h.svc.GetDoctorByID(clinicIDInt, req.DoctorID); err != nil {
        log.Println("Doctor not found:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
        return
    }

    schedule, err := h.svc.Create(clinicIDInt, &req)
    if err != nil {
        log.Println("Failed to create schedule:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, schedule)
}

// Update godoc
// @Summary Update a schedule
// @Description Update an existing schedule
// @Tags schedules
// @Accept json
// @Produce json
// @Param clinicID path int true "Clinic ID"
// @Param scheduleID path int true "Schedule ID"
// @Param schedule body dto.ScheduleRequest true "Schedule data"
// @Success 200 {object} models.Schedule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules/:scheduleID [put]
func (h *ScheduleHandler) Update(c *gin.Context) {
    log.Println("Received request to update schedule")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    scheduleID, err := strconv.Atoi(c.Param("scheduleID"))
    if err != nil {
        log.Println("Invalid schedule ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
        return
    }
    var req dto.ScheduleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    schedule, err := h.svc.Update(clinicIDInt, scheduleID, &req)
    if err != nil {
        log.Println("Failed to update schedule:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, schedule)
}

// Delete godoc
// @Summary Delete a schedule
// @Description Delete an existing schedule
// @Tags schedules
// @Accept json
// @Produce json
// @Param clinicID path int true "Clinic ID"
// @Param scheduleID path int true "Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /schedules/:scheduleID [delete]
func (h *ScheduleHandler) Delete(c *gin.Context) {
    log.Println("Received request to delete schedule")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    scheduleID, err := strconv.Atoi(c.Param("scheduleID"))
    if err != nil {
        log.Println("Invalid schedule ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
        return
    }
    if err := h.svc.Delete(clinicIDInt, scheduleID); err != nil {
        log.Println("Failed to delete schedule:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}