package handlers

import (
    "net/http"
    "vet-app-clinic-management/internal/service"
    "github.com/gin-gonic/gin"
	"strconv"
)

type DoctorHandler struct {
    svc *service.DoctorService
}

func NewDoctorHandler(svc *service.DoctorService) *DoctorHandler {
    return &DoctorHandler{svc: svc}
}

func (h *DoctorHandler) AuthenticateDoctor(c *gin.Context) {
    var req struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    doctor, err := h.svc.CreateOrGet(req.Name, req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Doctor authenticated", "doctor": doctor})
}

func (h *DoctorHandler) GetAll(c *gin.Context) {
    doctors, err := h.svc.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, doctors)
}

func (h *DoctorHandler) Update(c *gin.Context) {
    doctorID, err := strconv.Atoi(c.Param("doctorID"))
    if err != nil {
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

    doctor, err := h.svc.Update(doctorID, req.Name, req.Email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Doctor updated", "doctor": doctor})
}