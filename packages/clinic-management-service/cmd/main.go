package main

import (
    "context"
    "log"
    "net/http"
    "time"
    "clinic-management-service/internal/api"
    "clinic-management-service/internal/api/handlers"
    "clinic-management-service/internal/config"
    "clinic-management-service/internal/repository/mySQL"
    "clinic-management-service/internal/repository/redis"
    "clinic-management-service/internal/service"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Подключение к MySQL с повторными попытками
    var db *gorm.DB
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Failed to connect to MySQL, retrying in 5 seconds... (%v)", err)
        time.Sleep(5 * time.Second)
    }
    if err != nil {
        log.Fatalf("Failed to initialize MySQL after retries: %v", err)
    }

    // Подключение к Redis с повторными попытками
    var cache *redis.InventoryCache
    for i := 0; i < 10; i++ {
        cache = redis.NewInventoryCache(cfg.RedisAddr)
        ctx := context.Background()
        _, err := cache.Get(ctx, "test") // Проверяем подключение
        if err == nil || err.Error() == "redis: nil" { // "redis: nil" — это нормальный ответ, если ключа нет
            break
        }
        log.Printf("Failed to connect to Redis, retrying in 5 seconds... (%v)", err)
        time.Sleep(5 * time.Second)
    }
    if err != nil && err.Error() != "redis: nil" {
        log.Fatalf("Failed to initialize Redis after retries: %v", err)
    }

    // Репозитории
    scheduleRepo := mySQL.NewScheduleRepo(db)
    inventoryRepo := mySQL.NewInventoryRepo(db)

    // Сервисы
    scheduleSvc := service.NewScheduleService(scheduleRepo)
    inventorySvc := service.NewInventoryService(inventoryRepo, cache)

    // Хендлеры
    scheduleHandler := handlers.NewScheduleHandler(scheduleSvc)
    inventoryHandler := handlers.NewInventoryHandler(inventorySvc)

    // Роутер
    router := api.SetupRouter(scheduleHandler, inventoryHandler)

    // Запуск сервера
    log.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}