package main

import (
	"hospital-backend-api/config"
	"hospital-backend-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load()   // โหลดค่าจาก .env เข้า os.Getenv()

    config.ConnectDB() // เชื่อมต่อ PostgreSQL

    // gin.Default() สร้าง router พร้อม middleware พื้นฐาน
    r := gin.Default()
    routes.SetupRoutes(r) // ลงทะเบียน route ทั้งหมด

    r.Run(":8080") // เปิด server ที่ port 8080
}