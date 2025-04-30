package mySQL

import (
    "clinic-management-service/internal/models"
    "gorm.io/gorm"
)

type ScheduleRepo struct {
    db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
    return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) Create(schedule *models.Schedule) error {
    return r.db.Create(schedule).Error
}

func (r *ScheduleRepo) GetByDoctor(doctorID uint) ([]models.Schedule, error) {
    var schedules []models.Schedule
    err := r.db.Where("doctor_id = ?", doctorID).Find(&schedules).Error
    return schedules, err
}