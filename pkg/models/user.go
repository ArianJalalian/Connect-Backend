package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Age         uint
	Gender      string `gorm:"not null"`
	Password    string `gorm:"not null"`
	PhoneNumber string `gorm:"not null"`
	InfoID      string
	InfoType    string `gorm:"not null"`
	Block       bool
	Wallet      uint64
}
