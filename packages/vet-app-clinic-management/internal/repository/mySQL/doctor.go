package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
    "log"
)

type DoctorRepo struct {
    db *gorm.DB
}

func NewDoctorRepo(db *gorm.DB) *DoctorRepo {
    return &DoctorRepo{db: db}
}

func (r *DoctorRepo) GetByID(clinicID, id int) (*models.Doctor, error) {
    log.Printf("Fetching doctor by ID: %d for clinic ID %d", id, clinicID)
    var doctor models.Doctor
    if err := r.db.Where("clinic_id = ?", clinicID).First(&doctor, id).Error; err != nil {
        log.Println("Failed to fetch doctor by ID:", err)
        return nil, err
    }
    log.Printf("Doctor fetched by ID: ID=%d", doctor.ID)
    return &doctor, nil
}

func (r *DoctorRepo) GetByEmail(clinicID int, email string) (*models.Doctor, error) {
    log.Printf("Fetching doctor by email: %s for clinic ID %d", email, clinicID)
    var doctor models.Doctor
    if err := r.db.Where("email = ? AND clinic_id = ?", email, clinicID).First(&doctor).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            log.Println("Doctor not found by email")
            return nil, nil
        }
        log.Println("Failed to fetch doctor by email:", err)
        return nil, err
    }
    log.Printf("Doctor fetched by email: ID=%d", doctor.ID)
    return &doctor, nil
}

func (r *DoctorRepo) GetAll(clinicID int) ([]*models.Doctor, error) {
    log.Printf("Fetching all doctors for clinic ID %d from database", clinicID)
    var doctors []*models.Doctor
    if err := r.db.Where("clinic_id = ?", clinicID).Find(&doctors).Error; err != nil {
        log.Println("Failed to fetch all doctors from database:", err)
        return nil, err
    }
    log.Printf("Fetched %d doctors from database for clinic ID %d", len(doctors), clinicID)
    return doctors, nil
}

func (r *DoctorRepo) Create(doctor *models.Doctor) error {
    log.Printf("Creating doctor: name=%s, email=%s, clinic_id=%d", doctor.Name, doctor.Email, doctor.ClinicID)
    if err := r.db.Create(doctor).Error; err != nil {
        log.Println("Failed to create doctor in database:", err)
        return err
    }
    log.Printf("Doctor created in database: ID=%d", doctor.ID)
    return nil
}

func (r *DoctorRepo) Update(doctor *models.Doctor) error {
    log.Printf("Updating doctor in database: ID=%d", doctor.ID)
    if err := r.db.Save(doctor).Error; err != nil {
        log.Println("Failed to update doctor in database:", err)
        return err
    }
    log.Printf("Doctor updated in database: ID=%d", doctor.ID)
    return nil
}