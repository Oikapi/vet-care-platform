package db

import (
	"fmt"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/config"
	"github.com/Oikapi/vet-care-platform/packages/vet-app-appointment/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
	// Сначала подключаемся к PostgreSQL без указания конкретной базы данных
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)
	fmt.Println("Connecting to PostgreSQL with DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL:", err)
		return nil, err
	}

	// Проверяем, существует ли база данных, и создаём её, если нет
	dbName := cfg.DBName
	var dbExists bool
	err = db.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Scan(&dbExists).Error
	if err != nil {
		fmt.Println("Error checking if database exists:", err)
		return nil, err
	}

	if !dbExists {
		fmt.Println("Database does not exist, creating:", dbName)
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			fmt.Println("Failed to create database:", err)
			return nil, err
		}
	}

	// Подключаемся к созданной базе данных
	dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	fmt.Println("Connecting to PostgreSQL database with DSN:", dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL database:", err)
		return nil, err
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	// Автоматическая миграция схемы
	err = db.AutoMigrate(&models.Clinic{}, &models.Doctor{}, &models.Slot{}, &models.Appointment{})
	if err != nil {
		fmt.Println("Failed to auto-migrate schema:", err)
		return nil, err
	}
	fmt.Println("Schema migration completed")

	// Добавляем тестовые данные
	err = seedDatabase(db)
	if err != nil {
		fmt.Println("Failed to seed database:", err)
		return nil, err
	}
	fmt.Println("Database seeding completed")

	return db, nil
}

func seedDatabase(db *gorm.DB) error {
	// Проверяем, есть ли уже данные в таблице clinics
	var clinicCount int64
	if err := db.Model(&models.Clinic{}).Count(&clinicCount).Error; err != nil {
		return err
	}
	if clinicCount > 0 {
		fmt.Println("Database already seeded, skipping seeding")
		return nil
	}

	// Создаём тестовую клинику
	clinic := models.Clinic{
		Name:    "Vet Clinic A",
		Address: "123 Main St",
	}
	if err := db.Create(&clinic).Error; err != nil {
		return err
	}

	// Создаём тестового доктора
	doctor := models.Doctor{
		Name:      "Dr. Smith",
		Specialty: "Veterinarian",
		ClinicID:  clinic.ID,
	}
	if err := db.Create(&doctor).Error; err != nil {
		return err
	}

	// Создаём тестовые слоты
	slots := []models.Slot{
		{
			DoctorID: doctor.ID,
			SlotTime: time.Date(2025, 5, 2, 10, 0, 0, 0, time.UTC),
			IsBooked: false,
		},
		{
			DoctorID: doctor.ID,
			SlotTime: time.Date(2025, 5, 2, 11, 0, 0, 0, time.UTC),
			IsBooked: false,
		},
	}
	for _, slot := range slots {
		if err := db.Create(&slot).Error; err != nil {
			return err
		}
	}

	return nil
}

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return client, nil
}
