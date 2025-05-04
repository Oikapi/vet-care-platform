package mySQL

import (
    "vet-app-clinic-management/internal/models"
    "gorm.io/gorm"
)

type DoctorRepo struct {
    db *gorm.DB
}

func NewDoctorRepo(db *gorm.DB) *DoctorRepo {
    return &DoctorRepo{db: db}
}

func (r *DoctorRepo) GetByID(id int) (*models.Doctor, error) {
    var doctor models.Doctor
    if err := r.db.First(&doctor, id).Error; err != nil {
        return nil, err
    }
    return &doctor, nil
}

func (r *DoctorRepo) GetByEmail(email string) (*models.Doctor, error) {
    var doctor models.Doctor
    if err := r.db.Where("email = ?", email).First(&doctor).Error; err != nil {
        return nil, err
    }
    return &doctor, nil
}

func (r *DoctorRepo) GetAll() ([]*models.Doctor, error) {
    var doctors []*models.Doctor
    if err := r.db.Find(&doctors).Error; err != nil {
        return nil, err
    }
    return doctors, nil
}

func (r *DoctorRepo) Create(doctor *models.Doctor) error {
    return r.db.Create(doctor).Error
}

func (r *DoctorRepo) Update(doctor *models.Doctor) error {
    return r.db.Save(doctor).Error
}