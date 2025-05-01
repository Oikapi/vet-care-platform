package main

import (
	"fmt"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/config"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/db"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/handler"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/notification"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/repository"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/service"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация логгера
	log := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к PostgreSQL
	dbConn, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	// Подключение к Redis
	redisClient, err := db.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	// Инициализация репозиториев
	appointmentRepo := repository.NewAppointmentRepository(dbConn, redisClient)

	// Инициализация сервисов
	appointmentService := service.NewAppointmentService(appointmentRepo)

	// Инициализация Telegram Bot
	telegramBot, err := notification.NewTelegramBot(cfg.TelegramBotToken, appointmentService)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram Bot: %v", err)
	}

	// Запуск Telegram бота в отдельной горутине
	go telegramBot.Start()

	// Инициализация обработчиков
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, telegramBot)

	// Настройка Gin
	r := gin.Default()

	// Маршруты API
	api := r.Group("/api/v1")
	{
		api.POST("/appointments", appointmentHandler.CreateAppointment)
		api.GET("/appointments/slots/:clinic_id", appointmentHandler.GetAvailableSlots)
		api.GET("/appointments/:id", appointmentHandler.GetAppointment)
	}

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Infof("Сервер запущен на %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
