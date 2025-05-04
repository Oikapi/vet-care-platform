package service

import (
    "vet-app-clinic-management/internal/models"
    "vet-app-clinic-management/internal/repository/mySQL"
    "log"
)

type DoctorService struct {
    repo *mySQL.DoctorRepo
}

func NewDoctorService(repo *mySQL.DoctorRepo) *DoctorService {
    return &DoctorService{repo: repo}
}

func (s *DoctorService) GetByID(clinicID, id int) (*models.Doctor, error) {
    log.Printf("Fetching doctor with ID %d for clinic ID %d", id, clinicID)
    doctor, err := s.repo.GetByID(clinicID, id)
    if err != nil {
        log.Println("Failed to fetch doctor by ID:", err)
        return nil, err
    }
    log.Printf("Doctor fetched successfully: ID=%d", doctor.ID)
    return doctor, nil
}

func (s *DoctorService) GetByEmail(clinicID int, email string) (*models.Doctor, error) {
    log.Printf("Fetching doctor by email: %s for clinic ID %d", email, clinicID)
    doctor, err := s.repo.GetByEmail(clinicID, email)
    if err != nil {
        log.Println("Failed to fetch doctor by email:", err)
        return nil, err
    }
    if doctor == nil {
        log.Println("Doctor not found by email")
        return nil, nil
    }
    log.Printf("Doctor fetched by email: ID=%d", doctor.ID)
    return doctor, nil
}

func (s *DoctorService) GetAll(clinicID int) ([]*models.Doctor, error) {
    log.Printf("Fetching all doctors for clinic ID %d", clinicID)
    doctors, err := s.repo.GetAll(clinicID)
    if err != nil {
        log.Println("Failed to fetch all doctors:", err)
        return nil, err
    }
    log.Printf("Fetched %d doctors for clinic ID %d", len(doctors), clinicID)
    return doctors, nil
}

func (s *DoctorService) CreateOrGet(clinicID int, name, email string) (*models.Doctor, error) {
    log.Printf("Checking if doctor exists with email: %s for clinic ID %d", email, clinicID)
    doctor, err := s.repo.GetByEmail(clinicID, email)
    if err == nil && doctor != nil {
        log.Printf("Doctor found with email %s: ID=%d", email, doctor.ID)
        return doctor, nil
    }

    log.Printf("Doctor not found, creating new doctor: name=%s, email=%s for clinic ID %d", name, email, clinicID)
    doctor = &models.Doctor{
        Name:     name,
        Email:    email,
        ClinicID: clinicID,
    }
    if err := s.repo.Create(doctor); err != nil {
        log.Println("Failed to create doctor:", err)
        return nil, err
    }
    log.Printf("Doctor created successfully: ID=%d", doctor.ID)
    return doctor, nil
}

func (s *DoctorService) Update(clinicID, id int, name, email string) (*models.Doctor, error) {
    log.Printf("Fetching doctor with ID %d for update for clinic ID %d", id, clinicID)
    doctor, err := s.repo.GetByID(clinicID, id)
    if err != nil {
        log.Println("Failed to fetch doctor for update:", err)
        return nil, err
    }

    doctor.Name = name
    doctor.Email = email
    log.Printf("Updating doctor with ID %d: name=%s, email=%s", id, name, email)
    if err := s.repo.Update(doctor); err != nil {
        log.Println("Failed to update doctor:", err)
        return nil, err
    }
    log.Printf("Doctor updated successfully: ID=%d", doctor.ID)
    return doctor, nil
}