package models

type Inventory struct {
    ID           int    `json:"id" gorm:"primaryKey"`
    MedicineName string `json:"medicine_name"`
    Quantity     int    `json:"quantity"`
    Threshold    int    `json:"threshold"`
}