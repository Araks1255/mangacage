package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	AuthorID    uint
	TeamID      sql.NullInt32
	Team        Team    `gorm:"foreignKey:TeamID;references:id"`
	Genres      []Genre `gorm:"many2many:title_genres;"`
}
