package models

import (
	"time"

	"gorm.io/gorm"
)

type Clinic struct {
	gorm.Model
	ID       uint     `gorm:"primaryKey"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Photo    *string  `json:"photo,omitempty"`
	Doctors  []Doctor `gorm:"foreignKey:ClinicID" json:"doctors"`
}

type Doctor struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
	ClinicID       uint   `json:"clinicId"`
	Clinic         Clinic `json:"clinic" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Slot struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DoctorID  uint
	SlotTime  time.Time
	IsBooked  bool
}

type User struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Telegram  *string `json:"telegram,omitempty"`
	Password  string  `json:"password"`
	Photo     *string `json:"photo,omitempty"`
}

type Appointment struct {
	gorm.Model
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
