package main

import (
	"context"
	"fmt"
	"time"
	"vet-app-appointments/internal/client"
	"vet-app-appointments/internal/config"
	"vet-app-appointments/internal/handler"
	"vet-app-appointments/internal/notification"
	"vet-app-appointments/internal/repository"
	"vet-app-appointments/internal/service"
	"vet-app-appointments/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Инициализация логгера
	log := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к MySQL (vetcare_appointments)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var dbVetcare *gorm.DB
	for i := 0; i < 100; i++ {
		dbVetcare, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Infof("Попытка подключения к %s (%d/100): %v", cfg.DBName, i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to initialize MySQL (vetcare_appointments) after retries: %v", err)
	}

	// Подключение к MySQL (clinic_db)
	dsnClinic := fmt.Sprintf("%s:%s@tcp(%s:%s)/clinic_db?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	var dbClinic *gorm.DB
	for i := 0; i < 100; i++ {
		dbClinic, err = gorm.Open(mysql.Open(dsnClinic), &gorm.Config{})
		if err == nil {
			break
		}
		log.Infof("Попытка подключения к clinic_db (%d/100): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to initialize MySQL (clinic_db) after retries: %v", err)
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", cfg.RedisAddr),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Проверка соединения с Redis
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	// Инициализация репозиториев
	appointmentRepo := repository.NewAppointmentRepository(dbVetcare, dbClinic, redisClient)

	// Инициализация сервисов
	doctorClient := client.NewHTTPUserServiceClient("http://auth_backend:3000")
	appointmentService := service.NewAppointmentService(appointmentRepo, doctorClient)

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
	api := r.Group("appointments")
	{
		api.POST("/", appointmentHandler.CreateAppointment)
		api.GET("/slots/:clinic_id", appointmentHandler.GetAvailableSlots)
		api.GET("/:id", appointmentHandler.GetAppointment)
	}

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Infof("Сервер запущен на %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}