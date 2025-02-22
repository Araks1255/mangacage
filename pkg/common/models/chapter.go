package models

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	Name          string
	Description   string
	Path          string
	NumberOfPages int
	OnModeration  bool

	TitleID uint
	title   Title `gorm:"foreignKey:TitleID;references:id"`
}
