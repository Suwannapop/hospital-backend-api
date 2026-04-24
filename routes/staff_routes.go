package routes

import "github.com/gin-gonic/gin"

func setupStaffRoutes(r *gin.Engine) {
	api := r.Group("/staff")
		{
			api.GET("/", func(c *gin.Context) {
				c.String(200, "Hello, World!")
			})
		}
	
}