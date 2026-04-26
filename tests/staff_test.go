package tests

import (
	"fmt"
	"net/http"
	"testing"

	"hospital-backend-api/handler"
	"hospital-backend-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ==================== POST /staff/create ====================

func TestCreateStaff_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)

	hospital := seedHospital("Hospital A")

	body := map[string]interface{}{
		"username":    "nurse1",
		"password":    "pass123",
		"hospital_id": hospital.ID,
	}

	w := performRequest(r, "POST", "/staff/create", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(w)
	// password ต้องไม่อยู่ใน response
	fmt.Println("resp:", resp)
	assert.Equal(t, "nurse1", resp["username"])
	_, hasPassword := resp["password"]
	assert.False(t, hasPassword, "password ไม่ควรอยู่ใน response")
}

func TestCreateStaff_MissingFields(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)

	// ไม่ส่ง password
	body := map[string]interface{}{
		"username": "nurse1",
	}

	w := performRequest(r, "POST", "/staff/create", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseResponse(w)
	assert.Contains(t, resp["error"], "required")
}

func TestCreateStaff_DuplicateUsername(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)

	hospital := seedHospital("Hospital B")

	body := map[string]interface{}{
		"username":    "duplicate_user",
		"password":    "pass123",
		"hospital_id": hospital.ID,
	}

	// สร้างครั้งแรก — สำเร็จ
	w1 := performRequest(r, "POST", "/staff/create", body)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// สร้างครั้งที่สอง — ซ้ำ
	// หมายเหตุ: SQLite ส่ง "UNIQUE constraint failed" (500)
	// PostgreSQL ส่ง "duplicate key" ซึ่ง handler จับเป็น 409
	w2 := performRequest(r, "POST", "/staff/create", body)
	assert.True(t, w2.Code == http.StatusConflict || w2.Code == http.StatusInternalServerError,
		"ต้องเป็น 409 (PostgreSQL) หรือ 500 (SQLite)")
}

func TestCreateStaff_EmptyBody(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)

	w := performRequest(r, "POST", "/staff/create", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== POST /staff/login ====================

func TestLoginStaff_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)

	hospital := seedHospital("Hospital Login")

	// สร้าง staff ก่อน
	createBody := map[string]interface{}{
		"username":    "loginuser",
		"password":    "mypassword",
		"hospital_id": hospital.ID,
	}
	performRequest(r, "POST", "/staff/create", createBody)

	// Login
	loginBody := map[string]interface{}{
		"username":    "loginuser",
		"password":    "mypassword",
		"hospital_id": hospital.ID,
	}

	w := performRequest(r, "POST", "/staff/login", loginBody)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	assert.Equal(t, "เข้าสู่ระบบสำเร็จ", resp["message"])
	assert.NotEmpty(t, resp["token"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "loginuser", data["username"])
}

func TestLoginStaff_WrongPassword(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)

	hospital := seedHospital("Hospital WrongPW")

	createBody := map[string]interface{}{
		"username":    "user1",
		"password":    "correctpassword",
		"hospital_id": hospital.ID,
	}
	performRequest(r, "POST", "/staff/create", createBody)

	loginBody := map[string]interface{}{
		"username":    "user1",
		"password":    "wrongpassword",
		"hospital_id": hospital.ID,
	}

	w := performRequest(r, "POST", "/staff/login", loginBody)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := parseResponse(w)
	assert.Contains(t, resp["error"], "ไม่ถูกต้อง")
}

func TestLoginStaff_UserNotFound(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/login", handler.LoginStaff)

	seedHospital("Hospital NoUser")

	loginBody := map[string]interface{}{
		"username":    "nonexistent",
		"password":    "pass",
		"hospital_id": 1,
	}

	w := performRequest(r, "POST", "/staff/login", loginBody)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginStaff_MissingFields(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/login", handler.LoginStaff)

	// ส่งแค่ username
	loginBody := map[string]interface{}{
		"username": "user1",
	}

	w := performRequest(r, "POST", "/staff/login", loginBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginStaff_WrongHospital(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)

	hospital1 := seedHospital("Hospital X")
	hospital2 := seedHospital("Hospital Y")

	createBody := map[string]interface{}{
		"username":    "staffX",
		"password":    "pass123",
		"hospital_id": hospital1.ID,
	}
	performRequest(r, "POST", "/staff/create", createBody)

	// Login ด้วย hospital_id ผิด
	loginBody := map[string]interface{}{
		"username":    "staffX",
		"password":    "pass123",
		"hospital_id": hospital2.ID,
	}

	w := performRequest(r, "POST", "/staff/login", loginBody)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}


// ==================== Middleware: AuthRequired ====================

func TestAuthMiddleware_NoToken(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/protected", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := performRequest(r, "GET", "/protected", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := parseResponse(w)
	assert.Contains(t, resp["error"], "ไม่พบ token")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/protected", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := performRequestWithAuth(r, "GET", "/protected", "invalid-token-string", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)
	r.GET("/protected", middleware.AuthRequired(), func(c *gin.Context) {
		staffID, _ := c.Get("staff_id")
		hospitalID, _ := c.Get("hospital_id")
		c.JSON(http.StatusOK, gin.H{
			"staff_id":    staffID,
			"hospital_id": hospitalID,
		})
	})

	hospital := seedHospital("Hospital Auth")
	token := loginAndGetToken(r, "authuser", "pass123", hospital.ID)

	w := performRequestWithAuth(r, "GET", "/protected", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	assert.NotNil(t, resp["staff_id"])
	assert.NotNil(t, resp["hospital_id"])
}

