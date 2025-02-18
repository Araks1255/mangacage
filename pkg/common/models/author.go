package models

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name   string
	Genres []Genre `gorm:"many2many:author_genres;"`
}
