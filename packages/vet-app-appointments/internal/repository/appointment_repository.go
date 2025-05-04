package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vet-app-appointments/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

type AppointmentRepositoryInterface interface {
	CreateAppointment(appointment *models.Appointment) error
	GetAvailableSlots(clinicID uint, date time.Time) ([]models.Slot, error)
	GetAppointment(id uint) (*models.Appointment, error)
	BookSlot(slotID uint) error
	UnbookSlot(slotID uint) error
	CreateSlot(slot *models.Slot) error
	GetSlot(slotID uint) (*models.Slot, error)
	GetDoctor(doctorID uint) (*models.Doctor, error)
	GetClinic(clinicID uint) (*models.Clinic, error)
	UpdateAppointment(appointment *models.Appointment) error
}

func NewAppointmentRepository(db *gorm.DB, redis *redis.Client) *AppointmentRepository {
	return &AppointmentRepository{db: db, redis: redis}
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *AppointmentRepository) GetAvailableSlots(clinicID uint, date time.Time) ([]models.Slot, error) {
	cacheKey := fmt.Sprintf("slots:%d:%s", clinicID, date.Format("2006-01-02"))
	var slots []models.Slot

	// Проверяем кеш
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &slots); err == nil {
			return slots, nil
		}
	}

	// Запрашиваем из БД с JOIN на таблицу doctors
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	err = r.db.Joins("JOIN doctors ON doctors.id = slots.doctor_id").
		Where("doctors.clinic_id = ? AND slots.slot_time >= ? AND slots.slot_time < ? AND slots.is_booked = ?",
			clinicID, start, end, false).
		Find(&slots).Error
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	data, _ := json.Marshal(slots)
	r.redis.Set(context.Background(), cacheKey, data, 1*time.Hour)
	return slots, nil
}

func (r *AppointmentRepository) GetAppointment(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.First(&appointment, id).Error
	return &appointment, err
}

func (r *AppointmentRepository) BookSlot(slotID uint) error {
	return r.db.Model(&models.Slot{}).Where("id = ? AND is_booked = ?", slotID, false).
		Updates(map[string]interface{}{"is_booked": true}).Error
}

func (r *AppointmentRepository) UnbookSlot(slotID uint) error {
	return r.db.Model(&models.Slot{}).Where("id = ?", slotID).
		Updates(map[string]interface{}{"is_booked": false}).Error
}

func (r *AppointmentRepository) GetSlot(slotID uint) (*models.Slot, error) {
	var slot models.Slot
	err := r.db.First(&slot, slotID).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *AppointmentRepository) CreateSlot(slot *models.Slot) error {
	result := r.db.Create(slot)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AppointmentRepository) GetDoctor(doctorID uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.First(&doctor, doctorID).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *AppointmentRepository) GetClinic(clinicID uint) (*models.Clinic, error) {
	var clinic models.Clinic
	err := r.db.First(&clinic, clinicID).Error
	if err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (r *AppointmentRepository) UpdateAppointment(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}
