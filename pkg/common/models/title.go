package models

import (
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	AuthorID    uint
	TeamID      uint
	team        Team    `gorm:"foreignKey:TeamID;references:id"`
	genres      []Genre `gorm:"many2many:title_genres"`
}
