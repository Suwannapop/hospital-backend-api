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
	err := c.ShouldBindJSON(&hospital); 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
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