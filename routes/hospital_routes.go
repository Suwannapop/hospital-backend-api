package routes

import (
	"hospital-backend-api/handler"

	"github.com/gin-gonic/gin"
)

func setupHospitalRoutes(r *gin.Engine) {
	api := r.Group("/hospital")
	{
		api.POST("/", handler.CreateHospital) // เรียกใช้งานฟังก์ชันที่สร้างไว้
		api.GET("/", handler.GetHospitals)
	}
}