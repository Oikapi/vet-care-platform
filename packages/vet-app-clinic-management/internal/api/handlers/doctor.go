package handlers

import (
    "net/http"
    "vet-app-clinic-management/internal/service"
    "github.com/gin-gonic/gin"
    "strconv"
    "log"
)

type DoctorHandler struct {
    svc *service.DoctorService
}

func NewDoctorHandler(svc *service.DoctorService) *DoctorHandler {
    return &DoctorHandler{svc: svc}
}

func (h *DoctorHandler) AuthenticateDoctor(c *gin.Context) {
    log.Println("Received request to authenticate doctor")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    var req struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    doctor, err := h.svc.CreateOrGet(clinicIDInt, req.Name, req.Email)
    if err != nil {
        log.Println("Failed to authenticate doctor:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Doctor authenticated", "doctor": doctor})
}

func (h *DoctorHandler) GetAll(c *gin.Context) {
    log.Println("Received request to get all doctors")
    clinicID := c.Param("clinicID")
    id, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    doctors, err := h.svc.GetAll(id)
    if err != nil {
        log.Println("Failed to get all doctors:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, doctors)
}

func (h *DoctorHandler) Update(c *gin.Context) {
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    doctorID, err := strconv.Atoi(c.Param("doctorID"))
    if err != nil {
        log.Println("Invalid doctor ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
        return
    }

    var req struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    doctor, err := h.svc.Update(clinicIDInt, doctorID, req.Name, req.Email)
    if err != nil {
        log.Println("Failed to update doctor:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Doctor updated", "doctor": doctor})
}