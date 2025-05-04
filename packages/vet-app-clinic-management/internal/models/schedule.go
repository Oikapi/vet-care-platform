package models

import "time"

type Schedule struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    DoctorID  int       `json:"doctor_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Doctor    *Doctor   `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}