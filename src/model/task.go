package model

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Title       string `gorm:"size:255;not null;unique" json:"title"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Completed   bool   `gorm:"default:false" json:"completed"`
	UserId 			uint64
}
