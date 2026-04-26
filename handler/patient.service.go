package handler

import (
	"net/http"
	"strings"

	"hospital-backend-api/config"
	"hospital-backend-api/models"

	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	var patient models.Patient
	is_Error := c.ShouldBindJSON(&patient)
	if is_Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": is_Error.Error()})
		return
	}

	result := config.DB.Create(&patient)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "ข้อมูลซ้ำ: ไม่สามารถบันทึกได้เนื่องจากมีข้อมูลนี้อยู่ในระบบแล้ว"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	config.DB.Preload("Hospital").First(&patient, patient.ID) // 
	config.DB.Preload("Staff").First(&patient, patient.ID)
	c.JSON(http.StatusOK, struct {
		models.Patient
		Hospital any `json:"Hospital,omitempty"`
	}{
		Patient: patient,
		Hospital: struct {
			models.Hospital
			Patients any `json:"Patients,omitempty"` // ใช้ any แทน
		}{
			Hospital: patient.Hospital,
			Patients: nil,
		},
	})
}

func SearchPatientById(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
    
	// ค้นหาจาก NationalID หรือ PassportID ที่ตรงกับ id ที่ส่งเข้ามา
	result := config.DB.Where("national_id = ? OR passport_id = ?", id, id).First(&patient)
	
    if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ป่วย"})
		return
	}
	c.JSON(http.StatusOK, struct {
		FirstNameTH  string `json:"FirstNameTH"`
		MiddleNameTH string `json:"MiddleNameTH"`
		LastNameTH   string `json:"LastNameTH"`
		FirstNameEN  string `json:"FirstNameEN"`
		MiddleNameEN string `json:"MiddleNameEN"`
		LastNameEN   string `json:"LastNameEN"`
		DateOfBirth  string `json:"DateOfBirth"`
		PatientHN    string `json:"PatientHN"`
		NationalID   string `json:"NationalID"`
		PassportID   string `json:"PassportID"`
		PhoneNumber  string `json:"PhoneNumber"`
		Email        string `json:"Email"`
		Gender       string `json:"Gender"`
	}{
		FirstNameTH:  patient.FirstNameTH,
		MiddleNameTH: patient.MiddleNameTH,
		LastNameTH:   patient.LastNameTH,
		FirstNameEN:  patient.FirstNameEN,
		MiddleNameEN: patient.MiddleNameEN,
		LastNameEN:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth,
		PatientHN:    patient.PatientHN,
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       patient.Gender,
	})
}

// SearchPatient — ค้นหาผู้ป่วย (ต้อง login)
// แสดงเฉพาะผู้ป่วยที่อยู่ hospital เดียวกับ staff ที่ login
// Query params ทั้งหมด optional: national_id, passport_id, first_name, middle_name,
// last_name, date_of_birth, phone_number, email
func SearchPatient(c *gin.Context) {
	// ดึง hospital_id จาก JWT claims (middleware แนบไว้ใน context)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูล hospital จาก token"})
		return
	}

	// สร้าง Where Object เพื่อ serch
	query := config.DB.Where("hospital_id = ?", hospitalID)

	// เพิ่มเงื่อนไขตาม query params ที่ส่งมา (ทุก field optional)
	if v := c.Query("national_id"); v != "" {
		// ถ้ามี national_id and v
		query = query.Where("national_id = ?", v)
	}
	if v := c.Query("passport_id"); v != "" {
		query = query.Where("passport_id = ?", v)
	}
	if v := c.Query("first_name"); v != "" {
		query = query.Where("first_name_th ILIKE ? OR first_name_en ILIKE ?", "%"+v+"%", "%"+v+"%")
	}
	if v := c.Query("middle_name"); v != "" {
		query = query.Where("middle_name_th ILIKE ? OR middle_name_en ILIKE ?", "%"+v+"%", "%"+v+"%")
	}
	if v := c.Query("last_name"); v != "" {
		query = query.Where("last_name_th ILIKE ? OR last_name_en ILIKE ?", "%"+v+"%", "%"+v+"%")
	}
	if v := c.Query("date_of_birth"); v != "" {
		query = query.Where("date_of_birth = ?", v)
	}
	if v := c.Query("phone_number"); v != "" {
		query = query.Where("phone_number = ?", v)
	}
	if v := c.Query("email"); v != "" {
		query = query.Where("email ILIKE ?", "%"+v+"%")
	}

	var patients []models.Patient
	result := query.Find(&patients)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// แปลงผลลัพธ์เป็น response format
	type PatientResponse struct {
		NationalID   string `json:"national_id"`
		PassportID   string `json:"passport_id"`
		FirstNameTH  string `json:"first_name_th"`
		MiddleNameTH string `json:"middle_name_th"`
		LastNameTH   string `json:"last_name_th"`
		FirstNameEN  string `json:"first_name_en"`
		MiddleNameEN string `json:"middle_name_en"`
		LastNameEN   string `json:"last_name_en"`
		DateOfBirth  string `json:"date_of_birth"`
		PatientHN    string `json:"patient_hn"`
		PhoneNumber  string `json:"phone_number"`
		Email        string `json:"email"`
		Gender       string `json:"gender"`
	}

	var response []PatientResponse
	for _, p := range patients {
		response = append(response, PatientResponse{
			NationalID:   p.NationalID,
			PassportID:   p.PassportID,
			FirstNameTH:  p.FirstNameTH,
			MiddleNameTH: p.MiddleNameTH,
			LastNameTH:   p.LastNameTH,
			FirstNameEN:  p.FirstNameEN,
			MiddleNameEN: p.MiddleNameEN,
			LastNameEN:   p.LastNameEN,
			DateOfBirth:  p.DateOfBirth,
			PatientHN:    p.PatientHN,
			PhoneNumber:  p.PhoneNumber,
			Email:        p.Email,
			Gender:       p.Gender,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(response),
		"data":  response,
	})
}
