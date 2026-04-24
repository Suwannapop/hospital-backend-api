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
	err := c.ShouldBindJSON(&patient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Create(&patient)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			errMsg := "ข้อมูลซ้ำ: ไม่สามารถบันทึกได้เนื่องจากมีข้อมูลนี้อยู่ในระบบแล้ว"
			if strings.Contains(result.Error.Error(), "patient_hn") {
				errMsg = "ข้อมูลซ้ำ: รหัสผู้ป่วย (HN) นี้มีอยู่ในระบบแล้ว"
			} else if strings.Contains(result.Error.Error(), "national_id") {
				errMsg = "ข้อมูลซ้ำ: รหัสบัตรประชาชนนี้มีอยู่ในระบบแล้ว"
			} else if strings.Contains(result.Error.Error(), "passport_id") {
				errMsg = "ข้อมูลซ้ำ: รหัสพาสปอร์ตนี้มีอยู่ในระบบแล้ว"
			}
			
			c.JSON(http.StatusConflict, gin.H{
				"message": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	config.DB.Preload("Hospital").First(&patient, patient.ID)
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
	// c.JSON(http.StatusOK, patient)
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
	c.JSON(http.StatusOK, gin.H{
		"FirstNameTH":  patient.FirstNameTH,
		"MiddleNameTH": patient.MiddleNameTH,
		"LastNameTH":   patient.LastNameTH,
		"FirstNameEN":  patient.FirstNameEN,
		"MiddleNameEN": patient.MiddleNameEN,
		"LastNameEN":   patient.LastNameEN,
		"DateOfBirth":  patient.DateOfBirth,
		"PatientHN":    patient.PatientHN,
		"NationalID":   patient.NationalID,
		"PassportID":   patient.PassportID,
		"PhoneNumber":  patient.PhoneNumber,
		"Email":        patient.Email,
		"Gender":       patient.Gender,
	})
}
