package models

import "time"

type Hospital struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Staffs    []Staff   `gorm:"foreignKey:HospitalID"`
	Patients  []Patient `gorm:"foreignKey:HospitalID"`
}