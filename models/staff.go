package models

import "time"

type Staff struct {
    ID         uint      `gorm:"primaryKey"`
    HospitalID uint      `gorm:"not null;index"`
    Username   string    `gorm:"size:100;not null;unique"`
    Password   string    `gorm:"size:255;not null"`
    CreatedAt  time.Time
    UpdatedAt  time.Time

    Hospital   Hospital  `gorm:"foreignKey:HospitalID"`
}