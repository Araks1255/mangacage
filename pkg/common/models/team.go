package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `json:"name" binding:"required" gorm:"unique"`
	Description string `json:"description"`
}
