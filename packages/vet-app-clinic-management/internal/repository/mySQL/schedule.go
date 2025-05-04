package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
)

type ScheduleRepo struct {
    db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
    return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) GetByDoctorID(doctorID int) ([]*models.Schedule, error) {
    var schedules []*models.Schedule
    if err := r.db.Preload("Doctor").Where("doctor_id = ?", doctorID).Find(&schedules).Error; err != nil {
        return nil, err
    }
    return schedules, nil
}

func (r *ScheduleRepo) GetAll() ([]*models.Schedule, error) {
    var schedules []*models.Schedule
    if err := r.db.Preload("Doctor").Find(&schedules).Error; err != nil {
        return nil, err
    }
    return schedules, nil
}

func (r *ScheduleRepo) GetByID(id int) (*models.Schedule, error) {
    var schedule models.Schedule
    if err := r.db.Preload("Doctor").First(&schedule, id).Error; err != nil {
        return nil, err
    }
    return &schedule, nil
}

func (r *ScheduleRepo) Create(schedule *models.Schedule) error {
    return r.db.Create(schedule).Error
}

func (r *ScheduleRepo) Update(schedule *models.Schedule) error {
    return r.db.Save(schedule).Error
}

func (r *ScheduleRepo) Delete(scheduleID int) error {
    return r.db.Delete(&models.Schedule{}, scheduleID).Error
}