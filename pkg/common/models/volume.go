package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Volume struct {
	gorm.Model
	Name        string
	Description string

	TitleID uint
	Title   Title `gorm:"foreignKey:TitleID;references:id;OnDelete:SET NULL"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`
}
