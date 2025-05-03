package handler

import (
	"fmt"
	"vet-app-appointments/internal/models"
	"vet-app-appointments/internal/notification"
	"vet-app-appointments/internal/service"
	"vet-app-appointments/pkg/errors"
	"vet-app-appointments/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AppointmentHandler struct {
	service     *service.AppointmentService
	telegramBot *notification.TelegramBot
	log         *logger.Logger
}

func NewAppointmentHandler(service *service.AppointmentService, telegramBot *notification.TelegramBot) *AppointmentHandler {
	return &AppointmentHandler{
		service:     service,
		telegramBot: telegramBot,
		log:         logger.NewLogger(),
	}
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req struct {
		ClientID   uint   `json:"client_id" binding:"required"`
		DoctorID   uint   `json:"doctor_id" binding:"required"`
		ClinicID   uint   `json:"clinic_id" binding:"required"`
		SlotID     uint   `json:"slot_id" binding:"required"`
		TelegramID string `json:"telegram_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Неверный формат данных: %v", err)
		c.JSON(http.StatusBadRequest, errors.NewErrorResponse(fmt.Sprintf("Неверный формат данных: %v", err)))
		return
	}

	// Проверяем существование клиники
	_, err := h.service.GetClinic(req.ClinicID)
	if err != nil {
		h.log.Errorf("Клиника не найдена: %v", err)
		c.JSON(http.StatusBadRequest, errors.NewErrorResponse("Клиника не найдена"))
		return
	}

	// Проверяем существование доктора
	_, err = h.service.GetDoctor(req.DoctorID)
	if err != nil {
		h.log.Errorf("Доктор не найден: %v", err)
		c.JSON(http.StatusBadRequest, errors.NewErrorResponse("Доктор не найден"))
		return
	}

	// Проверяем существование слота
	slot, err := h.service.GetSlot(req.SlotID)
	if err != nil {
		h.log.Errorf("Слот не найден: %v", err)
		c.JSON(http.StatusBadRequest, errors.NewErrorResponse("Слот не найден"))
		return
	}

	appointment := &models.Appointment{
		ClientID:   req.ClientID,
		DoctorID:   req.DoctorID,
		ClinicID:   req.ClinicID,
		SlotID:     req.SlotID,
		TelegramID: req.TelegramID,
	}

	if err := h.service.CreateAppointment(appointment); err != nil {
		h.log.Errorf("Ошибка создания записи: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("Ошибка создания записи"))
		return
	}

	// Добавляем логирование ID созданной записи
	h.log.Infof("Запись успешно создана с ID: %d", appointment.ID)

	message := fmt.Sprintf("Напоминание: у вас запись в клинике на %s", slot.SlotTime.Format("2006-01-02 15:04"))
	if err := h.telegramBot.SendReminder(req.TelegramID, message); err != nil {
		h.log.Errorf("Ошибка отправки напоминания: %v", err)
		c.JSON(http.StatusOK, appointment)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	clinicID, _ := strconv.Atoi(c.Param("clinic_id"))
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.log.Errorf("Неверный формат даты: %v", err)
		c.JSON(http.StatusBadRequest, errors.NewErrorResponse("Неверный формат даты"))
		return
	}

	slots, err := h.service.GetAvailableSlots(uint(clinicID), date)
	if err != nil {
		h.log.Errorf("Ошибка получения слотов: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewErrorResponse("Ошибка получения слотов"))
		return
	}

	c.JSON(http.StatusOK, slots)
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	appointment, err := h.service.GetAppointment(uint(id))
	if err != nil {
		h.log.Errorf("Запись не найдена: %v", err)
		c.JSON(http.StatusNotFound, errors.NewErrorResponse("Запись не найдена"))
		return
	}

	c.JSON(http.StatusOK, appointment)
}
