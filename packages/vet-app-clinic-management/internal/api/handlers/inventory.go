package handlers

import (
    "vet-app-clinic-management/internal/api/dto"
    "vet-app-clinic-management/internal/service"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "log"
)

type InventoryHandler struct {
    svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
    return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) GetAll(c *gin.Context) {
    log.Println("Received request to get all inventory")
    clinicID := c.Param("clinicID")
    id, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    inventory, err := h.svc.GetAll(c.Request.Context(), id)
    if err != nil {
        log.Println("Failed to get all inventory:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) GetByID(c *gin.Context) {
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        log.Println("Invalid medicine ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
        return
    }

    inventory, err := h.svc.GetByID(c.Request.Context(), clinicIDInt, medicineID)
    if err != nil {
        log.Println("Failed to get inventory by ID:", err)
        if err.Error() == "record not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) AutoOrder(c *gin.Context) {
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    if err := h.svc.AutoOrder(c.Request.Context(), clinicIDInt); err != nil {
        log.Println("Failed to auto-order:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Auto-order completed"})
}

func (h *InventoryHandler) Create(c *gin.Context) {
    log.Println("Received request to create inventory")
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    var req dto.InventoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }
    inventory, err := h.svc.Create(c.Request.Context(), clinicIDInt, req.MedicineName, req.Quantity, req.Threshold)
    if err != nil {
        log.Println("Failed to create inventory:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, inventory)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        log.Println("Invalid medicine ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
        return
    }
    if err := h.svc.Delete(c.Request.Context(), clinicIDInt, medicineID); err != nil {
        log.Println("Failed to delete inventory:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted"})
}

func (h *InventoryHandler) UpdateQuantity(c *gin.Context) {
    clinicID := c.Param("clinicID")
    clinicIDInt, err := strconv.Atoi(clinicID)
    if err != nil {
        log.Println("Invalid clinic ID:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clinic ID"})
        return
    }
    medicineID, err := strconv.Atoi(c.Param("medicineID"))
    if err != nil {
        log.Println("Invalid medicine ID:", err)
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

    inventory, err := h.svc.UpdateQuantity(c.Request.Context(), clinicIDInt, medicineID, req.Quantity)
    if err != nil {
        log.Println("Failed to update quantity:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, inventory)
}