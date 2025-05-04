package service

import (
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "time"
)

type DoctorService struct {
    repo *mySQL.DoctorRepo
}

func NewDoctorService(repo *mySQL.DoctorRepo) *DoctorService {
    return &DoctorService{repo: repo}
}

func (s *DoctorService) GetByID(id int) (*models.Doctor, error) {
    return s.repo.GetByID(id)
}

func (s *DoctorService) GetByEmail(email string) (*models.Doctor, error) {
    return s.repo.GetByEmail(email)
}

func (s *DoctorService) GetAll() ([]*models.Doctor, error) {
    return s.repo.GetAll()
}

func (s *DoctorService) CreateOrGet(name, email string) (*models.Doctor, error) {
    doctor, err := s.GetByEmail(email)
    if err == nil {
        return doctor, nil
    }

    newDoctor := &models.Doctor{
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := s.repo.Create(newDoctor); err != nil {
        return nil, err
    }
    return newDoctor, nil
}

func (s *DoctorService) Update(id int, name, email string) (*models.Doctor, error) {
    doctor, err := s.GetByID(id)
    if err != nil {
        return nil, err
    }

    doctor.Name = name
    doctor.Email = email

    if err := s.repo.Update(doctor); err != nil {
        return nil, err
    }
    return doctor, nil
}