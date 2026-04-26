package handler

import (
	"net/http"

	"hospital-backend-api/config"
	"hospital-backend-api/models"

	"github.com/gin-gonic/gin"
)

func CreateHospital(c *gin.Context) {
    var hospital models.Hospital

    // 1. รับข้อมูลจาก request body
	is_Error := c.ShouldBindJSON(&hospital); 
    if is_Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": is_Error.Error(),
        })
        return
    }

    // 2. บันทึกลง database
	result := config.DB.Create(&hospital)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": result.Error.Error(),
        })
        return
    }

    // 3. ส่ง response กลับ
    c.JSON(http.StatusOK, hospital)
}

func GetHospitals(c *gin.Context) {
	var hospitals []models.Hospital
	result := config.DB.Find(&hospitals)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, hospitals)
}