package service

import (
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/models"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/repository"
	"time"
)

type AppointmentService struct {
	repo repository.AppointmentRepositoryInterface
}

func NewAppointmentService(repo repository.AppointmentRepositoryInterface) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) CreateAppointment(appointment *models.Appointment) error {
	appointment.Status = "confirmed"
	if err := s.repo.BookSlot(appointment.SlotID); err != nil {
		return err
	}
	return s.repo.CreateAppointment(appointment)
}

func (s *AppointmentService) GetAvailableSlots(clinicID uint, date time.Time) ([]models.Slot, error) {
	return s.repo.GetAvailableSlots(clinicID, date)
}

func (s *AppointmentService) GetAppointment(id uint) (*models.Appointment, error) {
	appointment, err := s.repo.GetAppointment(id)
	if err != nil {
		return nil, err
	}
	slot, err := s.repo.GetSlot(appointment.SlotID)
	if err != nil {
		return nil, err
	}
	appointment.Slot = *slot
	return appointment, nil
}

func (s *AppointmentService) GetSlot(slotID uint) (*models.Slot, error) {
	return s.repo.GetSlot(slotID)
}

func (s *AppointmentService) GetDoctor(doctorID uint) (*models.Doctor, error) {
	return s.repo.GetDoctor(doctorID)
}

func (s *AppointmentService) GetClinic(clinicID uint) (*models.Clinic, error) {
	return s.repo.GetClinic(clinicID)
}

func (s *AppointmentService) UpdateAppointment(appointment *models.Appointment) error {
	return s.repo.UpdateAppointment(appointment)
}

func (s *AppointmentService) UnbookSlot(slotID uint) error {
	return s.repo.UnbookSlot(slotID)
}
