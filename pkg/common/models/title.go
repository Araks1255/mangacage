package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Description  string
	OnModeration bool

	// Ещё обложку надо

	CreatorID uint
	User      User `gorm:"foreignKey:CreatorID;references:id"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id"`

	TeamID sql.NullInt32
	Team   Team `gorm:"foreignKey:TeamID;references:id"`

	Genres []Genre `gorm:"many2many:title_genres;"`
}
