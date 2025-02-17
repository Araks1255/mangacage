package models

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	Name        string
	Description string
	PathToFile  string

	TitleID uint
	title   Title `gorm:"foreignKey:TitleID;references:id"`
}
