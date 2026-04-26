package routes

import (
	"hospital-backend-api/handler"

	"github.com/gin-gonic/gin"
)

func setupStaffRoutes(r *gin.Engine) {
	api := r.Group("/staff")
		{
			api.POST("/create", handler.CreateStaff)
			api.POST("/login", handler.LoginStaff)
		}
	
}