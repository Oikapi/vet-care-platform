package api

import (
    "vet-app-clinic-management/internal/api/handlers"
    "github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты для API.
func SetupRouter(scheduleHandler *handlers.ScheduleHandler, inventoryHandler *handlers.InventoryHandler) *gin.Engine {
    router := gin.Default()

    apiGroup := router.Group("/management")
    {
        // Маршруты для расписаний
        apiGroup.GET("/schedules/doctor", scheduleHandler.GetAll)
        apiGroup.GET("/schedules/doctor/:doctorID", scheduleHandler.GetByDoctorID)
        apiGroup.POST("/schedules", scheduleHandler.Create)
        apiGroup.PUT("/schedules/:scheduleID", scheduleHandler.Update)
        apiGroup.DELETE("/schedules/:scheduleID", scheduleHandler.Delete)

        // Маршруты для инвентаря
        apiGroup.GET("/inventory", inventoryHandler.GetAll)
        apiGroup.GET("/inventory/:medicineID", inventoryHandler.GetByID)
        apiGroup.POST("/inventory/autoorder", inventoryHandler.AutoOrder)
        apiGroup.POST("/inventory", inventoryHandler.Create)
        apiGroup.DELETE("/inventory/:medicineID", inventoryHandler.Delete)
        apiGroup.PUT("/inventory/:medicineID/quantity", inventoryHandler.UpdateQuantity)
    }
    return router
}