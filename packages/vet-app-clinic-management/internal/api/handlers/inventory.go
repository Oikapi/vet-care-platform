package handlers

import (
	"net/http"
	"clinic-management-service/internal/service"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// ServeHTTP реализует интерфейс http.Handler
func (h *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAll godoc
// @Summary Get all inventory items
// @Description Retrieve complete list of inventory items
// @Tags inventory
// @Produce json
// @Success 200 {array} models.Inventory
// @Failure 500 {object} ErrorResponse
// @Router /api/inventory [get]
func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.svc.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, inventory)
}

