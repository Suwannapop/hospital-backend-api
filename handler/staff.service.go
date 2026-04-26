package handler

import (
	"net/http"
	"os"
	"strings"
	"time"

	"hospital-backend-api/config"
	"hospital-backend-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)



func CreateStaff(c *gin.Context) {
	var body struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		HospitalID uint   `json:"hospital_id" binding:"required"`
	}

	is_Error := c.ShouldBindJSON(&body)
	if is_Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": is_Error.Error()})
		return
	}

	// Hash password ด้วย bcrypt ก่อนเก็บลง DB
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถ hash รหัสผ่านได้"})
		return
	}

	staff := models.Staff{
		Username:   body.Username,
		Password:   string(hashedPassword),
		HospitalID: body.HospitalID,
	}

	result := config.DB.Create(&staff)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "username นี้มีอยู่ในระบบแล้ว"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Preload Hospital เพื่อส่งข้อมูล hospital กลับไปด้วย
	config.DB.Preload("Hospital").First(&staff, staff.ID)


	c.JSON(http.StatusCreated, struct {
		ID        uint           `json:"id"`
		Username  string         `json:"username"`
		Hospital  models.Hospital `json:"hospital"`
		CreatedAt time.Time      `json:"created_at"`
	}{
		ID:        staff.ID,
		Username:  staff.Username,
		Hospital:  staff.Hospital,
		CreatedAt: staff.CreatedAt,
	})
}

func LoginStaff(c *gin.Context) {
	var body struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		HospitalID uint   `json:"hospital_id" binding:"required"`
	}

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ค้นหา staff จาก username + hospital_id
	var staff models.Staff
	result := config.DB.Where("username = ? AND hospital_id = ?", body.Username, body.HospitalID).First(&staff)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username หรือ password ไม่ถูกต้อง"})
		return
	}

	// 3. เทียบ password ที่ส่งมากับ hash ใน DB
	err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username หรือ password ไม่ถูกต้อง"})
		return
	}

	// 4. สร้าง JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
		"username":    staff.Username,
		"exp":         time.Now().Add(24 * time.Hour).Unix(), // หมดอายุ 24 ชม.
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้าง token ได้"})
		return
	}

	// 5. ส่ง token กลับ
	c.JSON(http.StatusOK, gin.H{
		"message": "เข้าสู่ระบบสำเร็จ",
		"token":   tokenString,
		"data": gin.H{
			"staff_id":    staff.ID,
			"username":    staff.Username,
			"hospital_id": staff.HospitalID,
		},
	})
}

