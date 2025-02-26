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
	Creator   User `gorm:"foreignKey:CreatorID;references:id"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id"`

	TeamID sql.NullInt64
	Team   Team `gorm:"foreignKey:TeamID;references:id"`

	Genres []Genre `gorm:"many2many:title_genres;"`
}
