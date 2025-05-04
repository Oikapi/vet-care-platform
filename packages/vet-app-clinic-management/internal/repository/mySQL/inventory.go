package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
    "log"
)

type InventoryRepo struct {
    db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
    return &InventoryRepo{db: db}
}

func (r *InventoryRepo) GetAll(clinicID int) ([]*models.Inventory, error) {
    log.Printf("Fetching all inventory for clinic ID %d from database", clinicID)
    var inventories []*models.Inventory
    if err := r.db.Where("clinic_id = ?", clinicID).Find(&inventories).Error; err != nil {
        log.Println("Failed to fetch all inventory from database:", err)
        return nil, err
    }
    log.Printf("Fetched %d inventory items from database for clinic ID %d", len(inventories), clinicID)
    return inventories, nil
}

func (r *InventoryRepo) GetByID(clinicID, id int) (*models.Inventory, error) {
    log.Printf("Fetching inventory by ID: %d for clinic ID %d", id, clinicID)
    var inventory models.Inventory
    if err := r.db.Where("clinic_id = ?", clinicID).First(&inventory, id).Error; err != nil {
        log.Println("Failed to fetch inventory by ID:", err)
        return nil, err
    }
    log.Printf("Inventory fetched by ID: ID=%d", inventory.ID)
    return &inventory, nil
}

func (r *InventoryRepo) Update(inventory *models.Inventory) error {
    log.Printf("Updating inventory in database: ID=%d", inventory.ID)
    if err := r.db.Save(inventory).Error; err != nil {
        log.Println("Failed to update inventory in database:", err)
        return err
    }
    log.Printf("Inventory updated in database: ID=%d", inventory.ID)
    return nil
}

func (r *InventoryRepo) Create(inventory *models.Inventory) error {
    log.Printf("Creating inventory: medicine_name=%s, clinic_id=%d", inventory.MedicineName, inventory.ClinicID)
    if err := r.db.Create(inventory).Error; err != nil {
        log.Println("Failed to create inventory in database:", err)
        return err
    }
    log.Printf("Inventory created in database: ID=%d", inventory.ID)
    return nil
}

func (r *InventoryRepo) Delete(clinicID, id int) error {
    log.Printf("Deleting inventory with ID: %d for clinic ID %d from database", id, clinicID)
    if err := r.db.Where("clinic_id = ?", clinicID).Delete(&models.Inventory{}, id).Error; err != nil {
        log.Println("Failed to delete inventory from database:", err)
        return err
    }
    log.Printf("Inventory deleted from database: ID=%d", id)
    return nil
}