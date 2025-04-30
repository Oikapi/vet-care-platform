package models

import "time"

type Schedule struct {
    ID        uint      `gorm:"primaryKey"`
    DoctorID  uint      `gorm:"not null"`
    StartTime time.Time `gorm:"not null"`
    EndTime   time.Time `gorm:"not null"`
}