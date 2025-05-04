package dto

import "time"

type ScheduleRequest struct {
    DoctorID  int       `json:"doctor_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
}

type InventoryRequest struct {
    MedicineName string `json:"medicine_name"`
    Quantity     int    `json:"quantity"`
    Threshold    int    `json:"threshold"`
}