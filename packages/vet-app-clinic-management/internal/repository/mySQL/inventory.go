package mySQL

import (
    "context"
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
)

type InventoryRepo struct {
    db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
    return &InventoryRepo{db: db}
}

func (r *InventoryRepo) GetAll(ctx context.Context) ([]models.Inventory, error) {
    var inventory []models.Inventory
    err := r.db.Find(&inventory).Error
    return inventory, err
}

func (r *InventoryRepo) GetByID(id int) (*models.Inventory, error) {
    var inventory models.Inventory
    if err := r.db.First(&inventory, id).Error; err != nil {
        return nil, err
    }
    return &inventory, nil
}

func (r *InventoryRepo) Update(inventory *models.Inventory) error {
    return r.db.Save(inventory).Error
}

func (r *InventoryRepo) Create(inventory *models.Inventory) error {
    return r.db.Create(inventory).Error
}

func (r *InventoryRepo) Delete(medicineID int) error {
    return r.db.Delete(&models.Inventory{}, medicineID).Error
}