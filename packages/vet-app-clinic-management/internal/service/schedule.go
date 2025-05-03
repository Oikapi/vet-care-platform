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

func (s *ScheduleService) Create(schedule *models.Schedule) error {
    return s.repo.Create(schedule)
}

func (s *ScheduleService) GetByDoctor(doctorID uint) ([]models.Schedule, error) {
    return s.repo.GetByDoctor(doctorID)
}