package api

import (
    "vet-app-clinic-management/internal/api/handlers"
    "github.com/gin-gonic/gin"
)

func SetupRouter(scheduleHandler *handlers.ScheduleHandler, inventoryHandler *handlers.InventoryHandler, doctorHandler *handlers.DoctorHandler) *gin.Engine {
    router := gin.Default()

    apiGroup := router.Group("/management")
    {
        // Маршруты для расписаний
        apiGroup.GET(":clinicID/schedules/doctor", scheduleHandler.GetAll)
        apiGroup.GET(":clinicID/schedules/doctor/:doctorID", scheduleHandler.GetByDoctorID)
        apiGroup.POST(":clinicID/schedules", scheduleHandler.Create)
        apiGroup.PUT(":clinicID/schedules/:scheduleID", scheduleHandler.Update)
        apiGroup.DELETE(":clinicID/schedules/:scheduleID", scheduleHandler.Delete)

        // Маршруты для инвентаря
        apiGroup.GET(":clinicID/inventory", inventoryHandler.GetAll)
        apiGroup.GET(":clinicID/inventory/:medicineID", inventoryHandler.GetByID)
        apiGroup.POST(":clinicID/inventory/autoorder", inventoryHandler.AutoOrder)
        apiGroup.POST(":clinicID/inventory", inventoryHandler.Create)
        apiGroup.DELETE(":clinicID/inventory/:medicineID", inventoryHandler.Delete)
        apiGroup.PUT(":clinicID/inventory/:medicineID/quantity", inventoryHandler.UpdateQuantity)

        // Маршруты для врачей
        apiGroup.POST(":clinicID/auth/doctor", doctorHandler.AuthenticateDoctor)
        apiGroup.GET(":clinicID/doctors", doctorHandler.GetAll)
        apiGroup.PUT(":clinicID/doctors/:doctorID", doctorHandler.Update)
    }
    return router
}