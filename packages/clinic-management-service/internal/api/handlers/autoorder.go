package handlers

import "net/http"

func (h *InventoryHandler) AutoOrder(w http.ResponseWriter, r *http.Request) {
    if err := h.svc.CheckAndAutoOrder(r.Context()); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write([]byte("Auto-order completed"))
}