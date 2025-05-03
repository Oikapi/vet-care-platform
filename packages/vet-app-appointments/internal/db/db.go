package db

import (
	"fmt"
	"vet-app-appointments/internal/config"
	"vet-app-appointments/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL(cfg *config.Config) (*gorm.DB, error) {
	// Формируем DSN для подключения к MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	fmt.Println("Connecting to MySQL with DSN:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to MySQL:", err)
		return nil, err
	}
	fmt.Println("Successfully connected to MySQL database")

	// Автоматическая миграция схемы
	err = db.AutoMigrate(&models.Clinic{}, &models.Doctor{}, &models.Slot{}, &models.Appointment{})
	if err != nil {
		fmt.Println("Failed to auto-migrate schema:", err)
		return nil, err
	}
	fmt.Println("Schema migration completed")

	return db, nil
}

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return client, nil
}
