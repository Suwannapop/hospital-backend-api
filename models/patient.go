package models

import "time"

type Patient struct {
    ID           uint      `gorm:"primaryKey"`
    FirstNameTH  string    `gorm:"size:100"`
    MiddleNameTH string    `gorm:"size:100"`
    LastNameTH   string    `gorm:"size:100"`
    FirstNameEN  string    `gorm:"size:100"`
    MiddleNameEN string    `gorm:"size:100"`
    LastNameEN   string    `gorm:"size:100"`
    DateOfBirth  string    `gorm:"size:100"`
    PatientHN    string    `gorm:"size:50;index; unique"`
    NationalID   string    `gorm:"size:20;index; unique"`
    PassportID   string    `gorm:"size:20;index; unique"`
    PhoneNumber  string    `gorm:"size:20"`
    Email        string    `gorm:"size:100"`
    Gender       string    `gorm:"size:1"`

    HospitalID   uint      `gorm:"not null;index"`
    CreatedAt    time.Time
    UpdatedAt    time.Time

    Hospital     Hospital  `gorm:"foreignKey:HospitalID"`
}