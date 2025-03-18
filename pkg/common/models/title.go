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

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id;OnDelete:SET NULL"`

	TeamID sql.NullInt64
	Team   Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Genres []Genre `gorm:"many2many:title_genres;constraint:OnDelete:CASCADE"`
}
