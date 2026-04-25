package tests

import (
	"net/http"
	"testing"

	"hospital-backend-api/config"
	"hospital-backend-api/handler"
	"hospital-backend-api/middleware"
	"hospital-backend-api/models"

	"github.com/stretchr/testify/assert"
)

// ==================== POST /patient/create ====================

func TestCreatePatient_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/patient/create", handler.CreatePatient)

	// เรียก function seedHospital จาก test_helper.go เพื่อสร้างข้อมูลในฐานข้อมูล
	hospital := seedHospital("Hospital Patient")

	body := map[string]interface{}{
		"FirstNameTH": "สมชาย",
		"LastNameTH":  "ใจดี",
		"FirstNameEN": "Somchai",
		"LastNameEN":  "Jaidee",
		"DateOfBirth": "1990-01-01",
		"PatientHN":   "HN00001",
		"NationalID":  "1100000000001",
		"PassportID":  "AA123456",
		"PhoneNumber": "0812345678",
		"Email":       "somchai@example.com",
		"Gender":      "M",
		"HospitalID":  hospital.ID,
	}

	w := performRequest(r, "POST", "/patient/create", body)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePatient_EmptyBody(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/patient/create", handler.CreatePatient)

	w := performRequest(r, "POST", "/patient/create", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePatient_DuplicateNationalID(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/patient/create", handler.CreatePatient)

	hospital := seedHospital("Hospital Dup")

	body := map[string]interface{}{
		"FirstNameTH": "คนแรก",
		"LastNameTH":  "ทดสอบ",

		"PatientHN":   "HN00002",
		"NationalID":  "1100000000002",
		"PassportID":  "BB111111",

		"HospitalID":  hospital.ID,
	}

	w1 := performRequest(r, "POST", "/patient/create", body)
	assert.Equal(t, http.StatusOK, w1.Code)

	// สร้างซ้ำ (เปลี่ยน HN + Passport แต่ NationalID ซ้ำ)
	body["PatientHN"] = "HN00003"
	body["PassportID"] = "BB222222"

	w2 := performRequest(r, "POST", "/patient/create", body)
	// SQLite: 500 (UNIQUE constraint failed) / PostgreSQL: 409 (duplicate key)
	assert.True(t, w2.Code == http.StatusConflict || w2.Code == http.StatusInternalServerError,
		"ต้องเป็น 409 (PostgreSQL) หรือ 500 (SQLite)")
}

func TestCreatePatient_DuplicateHN(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/patient/create", handler.CreatePatient)

	hospital := seedHospital("Hospital DupHN")

	body := map[string]interface{}{
		"FirstNameTH": "คนแรก",
		"LastNameTH":  "ทดสอบ",
		"PatientHN":   "HN00002",
		"NationalID":  "1100000000002",
		"PassportID":  "BB111111",
		"HospitalID":  hospital.ID,
	}
	w1 := performRequest(r, "POST", "/patient/create", body)
	assert.Equal(t, http.StatusOK, w1.Code)

	body["NationalID"] = "1100000000003"
	body["PassportID"] = "BB222222"
	w2 := performRequest(r, "POST", "/patient/create", body)
	// SQLite: 500 (UNIQUE constraint failed) / PostgreSQL: 409 (duplicate key)
	assert.True(t, w2.Code == http.StatusConflict || w2.Code == http.StatusInternalServerError,
		"ต้องเป็น 409 (PostgreSQL) หรือ 500 (SQLite)")
}

// ==================== GET /patient/search/:id ====================

func TestSearchPatientBy_FoundByNationalID(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/patient/search/:id", handler.SearchPatientById)

	hospital := seedHospital("Hospital Search")

	patient := models.Patient{
		FirstNameTH: "สมหญิง",
		LastNameTH:  "รักดี",
		FirstNameEN: "Somying",
		LastNameEN:  "Rakdee",
		NationalID:  "1100000000099",
		PassportID:  "DD999999",
		PatientHN:   "HN00099",
		HospitalID:  hospital.ID,
	}
	config.DB.Create(&patient)

	// ค้นหาด้วย NationalID
	w := performRequest(r, "GET", "/patient/search/1100000000099", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	assert.Equal(t, "สมหญิง", resp["FirstNameTH"])
	assert.Equal(t, "รักดี", resp["LastNameTH"])
	assert.Equal(t, "Somying", resp["FirstNameEN"])
	assert.Equal(t, "Rakdee", resp["LastNameEN"])
	assert.Equal(t, "1100000000099", resp["NationalID"])
	assert.Equal(t, "DD999999", resp["PassportID"])
	assert.Equal(t, "HN00099", resp["PatientHN"])
}

func TestSearchPatientById_FoundByPassport(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/patient/search/:id", handler.SearchPatientById)

	hospital := seedHospital("Hospital Passport")

	patient := models.Patient{
		FirstNameEN: "John",
		LastNameEN:  "Doe",
		NationalID:  "1100000000088",
		PassportID:  "PP123456",
		PatientHN:   "HN00088",
		HospitalID:  hospital.ID,
	}
	config.DB.Create(&patient)

	// ค้นหาด้วย PassportID
	w := performRequest(r, "GET", "/patient/search/PP123456", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	assert.Equal(t, "John", resp["FirstNameEN"])
	assert.Equal(t, "Doe", resp["LastNameEN"])
	assert.Equal(t, "1100000000088", resp["NationalID"])
	assert.Equal(t, "PP123456", resp["PassportID"])
	assert.Equal(t, "HN00088", resp["PatientHN"])
}

func TestSearchPatientById_NotFound(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/patient/search/:id", handler.SearchPatientById)

	w := performRequest(r, "GET", "/patient/search/0000000000000", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := parseResponse(w)
	assert.Contains(t, resp["error"], "ไม่พบ")
}

// ==================== GET /patient/search (ต้อง login) ====================

func TestSearchPatient_WithLogin_ReturnsOwnHospital(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)
	r.GET("/patient/search", middleware.AuthRequired(), handler.SearchPatient)

	hospital1 := seedHospital("Hospital 1")
	hospital2 := seedHospital("Hospital 2")

	// สร้าง patient ใน hospital 1
	config.DB.Create(&models.Patient{
		FirstNameTH: "ผู้ป่วย1",
		NationalID:  "1111111111111",
		PassportID:  "P1111111",
		PatientHN:   "HN11111",
		HospitalID:  hospital1.ID,
	})

	// สร้าง patient ใน hospital 2
	config.DB.Create(&models.Patient{
		FirstNameTH: "ผู้ป่วย2",
		NationalID:  "2222222222222",
		PassportID:  "P2222222",
		PatientHN:   "HN22222",
		HospitalID:  hospital2.ID,
	})

	// login เป็น staff ของ hospital 1
	token := loginAndGetToken(r, "staff_h1", "pass123", hospital1.ID)

	// Search → ต้องเห็นเฉพาะ hospital 1
	w := performRequestWithAuth(r, "GET", "/patient/search", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	count := int(resp["count"].(float64))
	patient := resp["data"].([]any)[0].(map[string]any)
	assert.Equal(t, 1, count, "ต้องเห็นเฉพาะผู้ป่วยใน hospital ตัวเอง")
	assert.Equal(t, "ผู้ป่วย1", patient["first_name_th"])
	assert.Equal(t, "1111111111111", patient["national_id"])

}

func TestSearchPatient_WithFilter(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)
	r.GET("/patient/search", middleware.AuthRequired(), handler.SearchPatient)

	hospital := seedHospital("Hospital Filter")

	config.DB.Create(&models.Patient{
		FirstNameTH: "สมชาย",
		FirstNameEN: "Somchai",
		NationalID:  "3333333333333",
		PassportID:  "P3333333",
		PatientHN:   "HN33333",
		Email:       "somchai@test.com",
		HospitalID:  hospital.ID,
	})
	config.DB.Create(&models.Patient{
		FirstNameTH: "สมหญิง",
		FirstNameEN: "Somying",
		NationalID:  "4444444444444",
		PassportID:  "P4444444",
		PatientHN:   "HN44444",
		Email:       "somying@test.com",
		HospitalID:  hospital.ID,
	})

	token := loginAndGetToken(r, "staff_filter", "pass123", hospital.ID)

	// ค้นหาด้วย national_id
	w := performRequestWithAuth(r, "GET", "/patient/search?national_id=3333333333333", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	count := int(resp["count"].(float64))
	assert.Equal(t, 1, count)
	patient := resp["data"].([]any)[0].(map[string]any)
	assert.Equal(t, "สมชาย", patient["first_name_th"])
	assert.Equal(t, "3333333333333", patient["national_id"])
	assert.Equal(t, "somchai@test.com", patient["email"])
}

func TestSearchPatient_NoLogin_Returns401(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/patient/search", middleware.AuthRequired(), handler.SearchPatient)

	w := performRequest(r, "GET", "/patient/search", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSearchPatient_NoResults(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/staff/create", handler.CreateStaff)
	r.POST("/staff/login", handler.LoginStaff)
	r.GET("/patient/search", middleware.AuthRequired(), handler.SearchPatient)

	hospital := seedHospital("Hospital Empty")

	token := loginAndGetToken(r, "staff_empty", "pass123", hospital.ID)

	// ค้นหา — ไม่มี patient
	w := performRequestWithAuth(r, "GET", "/patient/search", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	count := int(resp["count"].(float64))
	assert.Equal(t, 0, count)
}
