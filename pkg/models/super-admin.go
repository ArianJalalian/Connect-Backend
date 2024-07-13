package models

import (
	"gorm.io/gorm"
)

type SuperAdmin struct {
	gorm.Model
	User   User
	UserID uint
}

type Report struct {
	gorm.Model
	Description string
	User        User
	UserID      uint
}
