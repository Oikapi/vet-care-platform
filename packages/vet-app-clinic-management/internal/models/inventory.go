package models

type Inventory struct {
    ID           uint   `gorm:"primaryKey"`
    MedicineName string `gorm:"not null"`
    Quantity     int    `gorm:"not null"`
    Threshold    int    `gorm:"not null"` // Порог для автозаказа
}