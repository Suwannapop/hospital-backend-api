package tests

import (
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
