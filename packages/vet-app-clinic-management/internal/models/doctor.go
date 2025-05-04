package models

import "time"

type Doctor struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
	ClinicID  int       `gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}