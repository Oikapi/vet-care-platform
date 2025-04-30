package handlers

import (
    "encoding/json"
    "net/http"
    "clinic-management-service/internal/service"
)

type InventoryHandler struct {
    svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
    return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    inventory, err := h.svc.GetAll(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(inventory)
}