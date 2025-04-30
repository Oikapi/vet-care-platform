package service

import (
    "clinic-management-service/internal/models"
    "clinic-management-service/internal/repository/mySQL"
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