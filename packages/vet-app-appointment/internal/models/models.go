package models

import (
	"gorm.io/gorm"
	"time"
)

type Clinic struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Address   string
}

type Doctor struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Specialty string
	ClinicID  uint
}

type Slot struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DoctorID  uint
	SlotTime  time.Time
	IsBooked  bool
}

type Appointment struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ClientID   uint
	DoctorID   uint
	ClinicID   uint
	SlotID     uint
	Slot       Slot `gorm:"foreignKey:SlotID"`
	Status     string
	TelegramID string
}
