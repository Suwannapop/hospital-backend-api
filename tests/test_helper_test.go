package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"hospital-backend-api/config"
	"hospital-backend-api/handler"
	"hospital-backend-api/middleware"
	"hospital-backend-api/models"
	"hospital-backend-api/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB สร้าง SQLite in-memory DB สำหรับ test
func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}

	db.AutoMigrate(&models.Hospital{}, &models.Staff{}, &models.Patient{})
	config.DB = db

	// ตั้ง JWT_SECRET สำหรับ test
	os.Setenv("JWT_SECRET", "test-secret-key")
}

// setupFullRouter สร้าง gin router พร้อม routes ทั้งหมด (เหมือน production)
func setupFullRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.SetupRoutes(r)
	return r
}

// setupRouter สร้าง gin router เปล่าสำหรับ test เฉพาะ handler
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// performRequest ส่ง request ไปที่ router แล้ว return response
func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request

	if body != nil {
		jsonBytes, _ := json.Marshal(body) //แปลง Go struct หรือ map → JSON bytes
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	//httptest.NewRecorder()สร้าง "กระดาษ" ไว้จด response (status, body, headers)
	// r.ServeHTTP(w, req) ยิง request เข้า router จริงๆ แล้วให้ router เขียน response ลงใน w
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// performRequestWithAuth ส่ง request พร้อม Authorization header
func performRequestWithAuth(r *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// seedHospital สร้าง hospital ใน test DB แล้ว return
func seedHospital(name string) models.Hospital {
	hospital := models.Hospital{Name: name}
	config.DB.Create(&hospital)
	return hospital
}

// parseResponse แปลง response body เป็น map
func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// loginAndGetToken สร้าง staff แล้ว login แล้ว return token
func loginAndGetToken(r *gin.Engine, username, password string, hospitalID uint) string {
	createBody := map[string]interface{}{
		"username":    username,
		"password":    password,
		"hospital_id": hospitalID,
	}
	performRequest(r, "POST", "/staff/create", createBody)

	loginW := performRequest(r, "POST", "/staff/login", createBody)
	resp := parseResponse(loginW)
	return resp["token"].(string)
}

// ใช้ import ทั้งหมดเพื่อไม่ให้ compiler บ่น
var (
	_ = handler.CreateStaff
	_ = middleware.AuthRequired
	_ = routes.SetupRoutes
)
