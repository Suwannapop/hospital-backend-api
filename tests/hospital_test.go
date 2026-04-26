package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"hospital-backend-api/handler"

	"github.com/stretchr/testify/assert"
)

// ==================== POST /hospital/ ====================

func TestCreateHospital_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/hospital/", handler.CreateHospital)

	body := map[string]string{
		"name": "โรงพยาบาลทดสอบ",
	}

	w := performRequest(r, "POST", "/hospital/", body)

	//assert.Equal(t,  expected,        actual  )
	assert.Equal(t, http.StatusOK, w.Code)
	// แปลง response body (JSON) → Go map เพื่อเอาค่าออกมาตรวจสอบ

	// fmt.Println("hold respond ",w)
	// fmt.Print("Body" , w.Body)
	resp := parseResponse(w)
	// fmt.Println(resp)
	assert.Equal(t, "โรงพยาบาลทดสอบ", resp["Name"])
	assert.NotNil(t, resp["ID"])
}
func TestCreateHospital_DuplicateName(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/hospital/", handler.CreateHospital)

	body := map[string]string{
		"name": "โรงพยาบาลซ้ำ",
	}

	// สร้างครั้งแรก — สำเร็จ
	w1 := performRequest(r, "POST", "/hospital/", body)
	assert.Equal(t, http.StatusOK, w1.Code)

	// สร้างครั้งที่สอง — ซ้ำ
	w2 := performRequest(r, "POST", "/hospital/", body)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
func TestCreateHospital_EmptyBody(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/hospital/", handler.CreateHospital)

	w := performRequest(r, "POST", "/hospital/", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHospital_InvalidJSON(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.POST("/hospital/", handler.CreateHospital)

	
	// ส่ง string ที่ไม่ใช่ JSON
	w := performRequest(r, "POST", "/hospital/", "invalid")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHospital_WithAddionalField(t *testing.T){
	setupTestDB()
	r := setupRouter()
	r.POST("/hospital/", handler.CreateHospital)

	body := map[string]interface{}{
		"name": "โรงพยาบาลทดสอบ",
		"description": "description",
	}
	
	w := performRequest(r, "POST", "/hospital/", body)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(w)
	assert.Equal(t, "โรงพยาบาลทดสอบ", resp["Name"])
	assert.NotNil(t, resp["ID"])
	assert.Nil(t , resp["description"])

}

// ==================== GET /hospital/ ====================

func TestGetHospitals_Success(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/hospital/", handler.GetHospitals)

	// สร้าง hospital ใน DB ก่อน
	seedHospital("โรงพยาบาล A")
	seedHospital("โรงพยาบาล B")

	w := performRequest(r, "GET", "/hospital/", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	// parse เป็น array
	var hospitals []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &hospitals)

	assert.Equal(t, 2, len(hospitals))
	assert.Equal(t, "โรงพยาบาล A", hospitals[0]["Name"])
	assert.Equal(t, "โรงพยาบาล B", hospitals[1]["Name"])
}

func TestGetHospitals_EmptyList(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/hospital/", handler.GetHospitals)

	// ไม่ seed ข้อมูล — DB ว่าง
	w := performRequest(r, "GET", "/hospital/", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var hospitals []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &hospitals)

	// ต้องได้ array ว่าง หรือ null
	assert.True(t, len(hospitals) == 0 || hospitals == nil)
}

func TestGetHospitals_ReturnsAllFields(t *testing.T) {
	setupTestDB()
	r := setupRouter()
	r.GET("/hospital/", handler.GetHospitals)

	seedHospital("โรงพยาบาลตรวจ Fields")

	w := performRequest(r, "GET", "/hospital/", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var hospitals []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &hospitals)

	assert.Equal(t, 1, len(hospitals))
	h := hospitals[0]
	assert.NotNil(t, h["ID"])
	assert.Equal(t, "โรงพยาบาลตรวจ Fields", h["Name"])
	assert.NotNil(t, h["CreatedAt"])
	assert.NotNil(t, h["UpdatedAt"])
}

