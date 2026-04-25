package routes

import (
	"hospital-backend-api/handler"

	"github.com/gin-gonic/gin"
)

func setupStaffRoutes(r *gin.Engine) {
	api := r.Group("/staff")
		{
			api.GET("/", func(c *gin.Context) {
				c.String(200, "Hello, World!")
			})
			api.POST("/create", handler.CreateStaff)
			api.POST("/login", handler.LoginStaff)
			api.POST("/logout", handler.LogoutStaff)
		}
	
}