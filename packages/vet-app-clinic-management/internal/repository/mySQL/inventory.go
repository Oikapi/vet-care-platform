package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
)

type InventoryRepo struct {
    db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
    return &InventoryRepo{db: db}
}

func (r *InventoryRepo) GetAll() ([]models.Inventory, error) {
    var inventory []models.Inventory
    err := r.db.Find(&inventory).Error
    return inventory, err
}

func (r *InventoryRepo) UpdateQuantity(id uint, quantity int) error {
    return r.db.Model(&models.Inventory{}).Where("id = ?", id).Update("quantity", quantity).Error
}