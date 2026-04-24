package handler

import (
	"net/http"

	"hospital-backend-api/models"

	"github.com/gin-gonic/gin"
)

func CreateStaff(c *gin.Context) {
	var staff models.Staff
	err := c.ShouldBindJSON(&staff)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
