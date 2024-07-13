package models

import (
	"gorm.io/gorm"
)

type Trainee struct {
	gorm.Model
	User           User
	UserID         uint
	Height         uint
	Weight         uint
	Sports         string
	UserName       string
	Status         string
	Contact        string
	Language       string
	Country        string
	MedicalHistory string
	ActiveDays     ActiveDays
	ActiveDaysID   uint
	RequestID      uint `gorm:"column:request_id"`
	TrainerID      uint
}
