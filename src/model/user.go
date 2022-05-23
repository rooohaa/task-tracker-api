package model

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Verified bool   `gorm:"default:false" json:"is_verified"`
}
