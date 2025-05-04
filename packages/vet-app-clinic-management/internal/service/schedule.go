package service

import (
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
)

type ScheduleService struct {
    repo *mySQL.ScheduleRepo
}

func NewScheduleService(repo *mySQL.ScheduleRepo) *ScheduleService {
    return &ScheduleService{repo: repo}
}

func (s *ScheduleService) GetAll() ([]*models.Schedule, error) {
    return s.repo.GetAll()
}

func (s *ScheduleService) GetByDoctorID(doctorID int) ([]models.Schedule, error) {
    return s.repo.GetByDoctorID(doctorID)
}

func (s *ScheduleService) Create(schedule *models.Schedule) error {
    return s.repo.Create(schedule)
}

func (s *ScheduleService) Update(schedule *models.Schedule) error {
    return s.repo.Update(schedule)
}

func (s *ScheduleService) Delete(scheduleID int) error {
    return s.repo.Delete(scheduleID)
}