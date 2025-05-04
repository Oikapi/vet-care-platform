package handlers

import (
    "vet-app-clinic-management/internal/api/dto"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/service"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

type InventoryHandler struct {
    svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
    return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) GetAll(c *gin.Context) {
    inventory, err := h.svc.GetAll(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) GetByID(c *gin.Context) {
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
        return
    }

    inventory, err := h.svc.GetByID(c.Request.Context(), medicineID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
        return
    }
    c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) AutoOrder(c *gin.Context) {
    if err := h.svc.AutoOrder(c.Request.Context()); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Auto-order completed"})
}

func (h *InventoryHandler) Create(c *gin.Context) {
    var req dto.InventoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    inventory := &models.Inventory{
        MedicineName: req.MedicineName,
        Quantity:     req.Quantity,
        Threshold:    req.Threshold,
    }
    if err := h.svc.Create(c.Request.Context(), inventory); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, inventory)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
        return
    }
    if err := h.svc.Delete(c.Request.Context(), medicineID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted"})
}

func (h *InventoryHandler) UpdateQuantity(c *gin.Context) {
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
        return
    }

    var req struct {
        Quantity int `json:"Quantity" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    inventory, err := h.svc.UpdateQuantity(c.Request.Context(), medicineID, req.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, inventory)
}