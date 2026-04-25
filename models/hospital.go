package models

import "time"

type Hospital struct {
	ID        uint      `gorm:"primaryKey"`

	// ไม่มี json:"..." tag → Go ใช้ชื่อ field (Name) แล้วจับคู่แบบ ไม่สนตัวเล็กตัวใหญ่
	Name      string    `gorm:"not null;unique"`

	// Name string `gorm:"not null;unique" json:"name"` 
	// ถ้าอยากบังคับให้ใช้แค่ name (ตัวเล็ก) เพิ่ม json tag ที่ model:


	CreatedAt time.Time
	UpdatedAt time.Time

	Staffs    []Staff   `gorm:"foreignKey:HospitalID"`
	Patients  []Patient `gorm:"foreignKey:HospitalID"`
}