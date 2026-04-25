package config

import (
	"fmt"
	"log"
	"os"

	"hospital-backend-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {

    // สร้าง connection string
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )


    //เชื่อมต่อ DB ผ่าน GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        // print SQL ทุก query (debug ง่ายมาก)
        // ตอน production เปลี่ยนเป็น logger.Silent
    })

    if err != nil {
        log.Fatal("❌ เชื่อมต่อ DB ไม่ได้:", err)
    }

    // Connection Pool — จัดการ connection ให้มีประสิทธิภาพ
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("❌ ดึง sql.DB ไม่ได้:", err)
    }
    sqlDB.SetMaxOpenConns(10)   // connection พร้อมกันสูงสุด 10 อัน
    sqlDB.SetMaxIdleConns(5)    // connection ไว้เฉย ๆ” ได้สูงสุด 5 อัน
    
    DB = db
    log.Println("✅ เชื่อมต่อ DB สำเร็จ")

    // Auto-migrate — สร้าง/อัปเดต table ให้ตรงกับ struct
    // ลำดับสำคัญ Hospital ต้องก่อน เพราะ Staff และ Patient อ้างอิงถึง
    err = DB.AutoMigrate(
        &models.Hospital{},
        &models.Staff{},
        &models.Patient{},
    )
    if err != nil {
        log.Fatal("❌ Auto-migrate ไม่ได้:", err)
    }
    log.Println("✅ Auto-migrate สำเร็จ")
}