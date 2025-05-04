package main

import (
    "context"
    "log"
    "time"
    "vet-app-clinic-management/internal/api"
    "vet-app-clinic-management/internal/api/handlers"
    "vet-app-clinic-management/internal/config"
    "vet-app-clinic-management/internal/repository/mySQL"
    "vet-app-clinic-management/internal/repository/redis"
    "vet-app-clinic-management/internal/service"
    "vet-app-clinic-management/internal/docs"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/swaggo/http-swagger"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Подключение к MySQL
    var db *gorm.DB
    for i := 0; i < 100; i++ {
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

    // Подключение к Redis
    var cache *redis.InventoryCache
    for i := 0; i < 10; i++ {
        cache, err = redis.NewInventoryCache(cfg.RedisAddr)
        if err != nil {
            log.Printf("Failed to connect to Redis, retrying in 5 seconds... (%v)", err)
            time.Sleep(5 * time.Second)
            continue
        }
        ctx := context.Background()
        _, err := cache.Get(ctx, "test")
        if err == nil || err.Error() == "redis: nil" {
            break
        }
        log.Printf("Failed to connect to Redis, retrying in 5 seconds... (%v)", err)
        time.Sleep(5 * time.Second)
    }
    if err != nil && err.Error() != "redis: nil" {
        log.Fatalf("Failed to initialize Redis after retries: %v", err)
    }

    // Инициализация сервисов и обработчиков
    scheduleRepo := mySQL.NewScheduleRepo(db)
    inventoryRepo := mySQL.NewInventoryRepo(db)
    scheduleSvc := service.NewScheduleService(scheduleRepo)
    inventorySvc := service.NewInventoryService(inventoryRepo, cache)
    scheduleHandler := handlers.NewScheduleHandler(scheduleSvc)
    inventoryHandler := handlers.NewInventoryHandler(inventorySvc)

    // Настройка Swagger
    docs.SwaggerInfo.Title = "Clinic Management API"
    docs.SwaggerInfo.Description = "API for clinic management"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = cfg.SwaggerHost
    docs.SwaggerInfo.BasePath = "/management"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

    // Настройка роутера с gin
    router := api.SetupRouter(scheduleHandler, inventoryHandler)

    // Добавляем Swagger UI
    router.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))
    // Добавляем raw swagger.json
    router.GET("/swagger.json", func(c *gin.Context) {
        file, _ := docs.SwaggerJSON.ReadFile("swagger.json")
        c.Header("Content-Type", "application/json")
        c.Writer.Write(file)
    })

    // Запуск сервера
    log.Printf("Starting server on :8080, Swagger UI at http://%s/swagger/index.html", cfg.SwaggerHost)
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}