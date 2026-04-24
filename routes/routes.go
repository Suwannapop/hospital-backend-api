package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	setupStaffRoutes(r)
	setupHospitalRoutes(r)
	setupPatientRoutes(r)
}