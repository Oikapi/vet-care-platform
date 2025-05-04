package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
    "log"
)

type ScheduleRepo struct {
    db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
    return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) GetByDoctorID(clinicID, doctorID int) ([]*models.Schedule, error) {
    log.Printf("Fetching schedules for doctor ID: %d and clinic ID %d from database", doctorID, clinicID)
    var schedules []*models.Schedule
    if err := r.db.Preload("Doctor").Where("doctor_id = ? AND clinic_id = ?", doctorID, clinicID).Find(&schedules).Error; err != nil {
        log.Println("Failed to fetch schedules by doctor ID:", err)
        return nil, err
    }
    log.Printf("Fetched %d schedules for doctor ID %d and clinic ID %d", len(schedules), doctorID, clinicID)
    return schedules, nil
}

func (r *ScheduleRepo) GetAll(clinicID int) ([]*models.Schedule, error) {
    log.Printf("Fetching all schedules for clinic ID %d from database with doctor details", clinicID)
    var schedules []*models.Schedule
    if err := r.db.Preload("Doctor").Where("clinic_id = ?", clinicID).Find(&schedules).Error; err != nil {
        log.Println("Failed to fetch all schedules from database:", err)
        return nil, err
    }
    log.Printf("Fetched %d schedules from database for clinic ID %d", len(schedules), clinicID)
    return schedules, nil
}

func (r *ScheduleRepo) GetByID(clinicID, id int) (*models.Schedule, error) {
    log.Printf("Fetching schedule by ID: %d for clinic ID %d from database", id, clinicID)
    var schedule models.Schedule
    if err := r.db.Preload("Doctor").Where("clinic_id = ?", clinicID).First(&schedule, id).Error; err != nil {
        log.Println("Failed to fetch schedule by ID:", err)
        return nil, err
    }
    log.Printf("Schedule fetched by ID: ID=%d", schedule.ID)
    return &schedule, nil
}

func (r *ScheduleRepo) Create(schedule *models.Schedule) error {
    log.Printf("Creating schedule: doctor_id=%d, clinic_id=%d, start_time=%v, end_time=%v", schedule.DoctorID, schedule.ClinicID, schedule.StartTime, schedule.EndTime)
    if err := r.db.Create(schedule).Error; err != nil {
        log.Println("Failed to create schedule in database:", err)
        return err
    }
    log.Printf("Schedule created in database: ID=%d", schedule.ID)
    return nil
}

func (r *ScheduleRepo) Update(schedule *models.Schedule) error {
    log.Printf("Updating schedule in database: ID=%d", schedule.ID)
    if err := r.db.Save(schedule).Error; err != nil {
        log.Println("Failed to update schedule in database:", err)
        return err
    }
    log.Printf("Schedule updated in database: ID=%d", schedule.ID)
    return nil
}

func (r *ScheduleRepo) Delete(clinicID, id int) error {
    log.Printf("Deleting schedule with ID: %d for clinic ID %d from database", id, clinicID)
    if err := r.db.Where("clinic_id = ?", clinicID).Delete(&models.Schedule{}, id).Error; err != nil {
        log.Println("Failed to delete schedule from database:", err)
        return err
    }
    log.Printf("Schedule deleted from database: ID=%d", id)
    return nil
}