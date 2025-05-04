package service

import (
    "vet-app-clinic-management/internal/api/dto"
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "log"
    "errors"
)

type ScheduleService struct {
    repo       *mySQL.ScheduleRepo
    doctorRepo *mySQL.DoctorRepo
}

func NewScheduleService(repo *mySQL.ScheduleRepo, doctorRepo *mySQL.DoctorRepo) *ScheduleService {
    return &ScheduleService{repo: repo, doctorRepo: doctorRepo}
}

func (s *ScheduleService) GetAll(clinicID int) ([]*models.Schedule, error) {
    log.Printf("Fetching all schedules with doctor details for clinic ID %d", clinicID)
    schedules, err := s.repo.GetAll(clinicID)
    if err != nil {
        log.Println("Failed to fetch all schedules:", err)
        return nil, err
    }
    log.Printf("Fetched %d schedules for clinic ID %d", len(schedules), clinicID)
    return schedules, nil
}

func (s *ScheduleService) GetByID(clinicID, scheduleID int) (*models.Schedule, error) {
    log.Printf("Fetching schedule with ID %d for clinic ID %d", scheduleID, clinicID)
    schedule, err := s.repo.GetByID(clinicID, scheduleID)
    if err != nil {
        log.Println("Failed to fetch schedule by ID:", err)
        return nil, err
    }
    log.Printf("Schedule fetched successfully: ID=%d", schedule.ID)
    return schedule, nil
}

func (s *ScheduleService) GetByDoctorID(clinicID, doctorID int) ([]*models.Schedule, error) {
    log.Printf("Fetching schedules for doctor ID: %d and clinic ID %d", doctorID, clinicID)
    schedules, err := s.repo.GetByDoctorID(clinicID, doctorID)
    if err != nil {
        log.Println("Failed to fetch schedules by doctor ID:", err)
        return nil, err
    }
    log.Printf("Fetched %d schedules for doctor ID %d and clinic ID %d", len(schedules), doctorID, clinicID)
    return schedules, nil
}

func (s *ScheduleService) Create(clinicID int, req *dto.ScheduleRequest) (*models.Schedule, error) {
    log.Printf("Creating schedule for doctor ID: %d and clinic ID %d", req.DoctorID, clinicID)
    doctor, err := s.doctorRepo.GetByID(clinicID, req.DoctorID)
    if err != nil {
        log.Println("Failed to fetch doctor:", err)
        return nil, err
    }
    if doctor == nil {
        log.Println("Doctor not found")
        return nil, errors.New("doctor not found")
    }
    schedule := &models.Schedule{
        DoctorID:  req.DoctorID,
        ClinicID:  clinicID,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
    }
    if err := s.repo.Create(schedule); err != nil {
        log.Println("Failed to create schedule:", err)
        return nil, err
    }
    log.Printf("Schedule created successfully: ID=%d", schedule.ID)
    return schedule, nil
}

func (s *ScheduleService) Update(clinicID, id int, req *dto.ScheduleRequest) (*models.Schedule, error) {
    log.Printf("Updating schedule with ID: %d for clinic ID %d", id, clinicID)
    schedule, err := s.repo.GetByID(clinicID, id)
    if err != nil {
        log.Println("Failed to fetch schedule for update:", err)
        return nil, err
    }
    schedule.StartTime = req.StartTime
    schedule.EndTime = req.EndTime
    if err := s.repo.Update(schedule); err != nil {
        log.Println("Failed to update schedule:", err)
        return nil, err
    }
    log.Printf("Schedule updated successfully: ID=%d", schedule.ID)
    return schedule, nil
}

func (s *ScheduleService) Delete(clinicID, id int) error {
    log.Printf("Deleting schedule with ID: %d for clinic ID %d", id, clinicID)
    if err := s.repo.Delete(clinicID, id); err != nil {
        log.Println("Failed to delete schedule:", err)
        return err
    }
    log.Printf("Schedule deleted successfully: ID=%d", id)
    return nil
}

func (s *ScheduleService) GetDoctorByID(clinicID, doctorID int) (*models.Doctor, error) {
    log.Printf("Fetching doctor with ID %d for clinic ID %d", doctorID, clinicID)
    doctor, err := s.doctorRepo.GetByID(clinicID, doctorID)
    if err != nil {
        log.Println("Failed to fetch doctor by ID:", err)
        return nil, err
    }
    log.Printf("Doctor fetched successfully: ID=%d", doctor.ID)
    return doctor, nil
}