package routes

import (
	"hospital-backend-api/handler"

	"github.com/gin-gonic/gin"
)

func setupPatientRoutes(r *gin.Engine) {
	api := r.Group("/patient")
	{
		api.GET("/", func(c *gin.Context) {
			c.String(200, "Hello, World!")
		})
		api.POST("/create", handler.CreatePatient)
		api.GET("/search/:id", handler.SearchPatientById)
	}
}