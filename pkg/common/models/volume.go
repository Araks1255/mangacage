package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Volume struct {
	gorm.Model

	Name         string
	Description  string
	OnModeration bool

	TitleID uint
	Title   Title `gorm:"foreignKey:TitleID;references:id"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id"`
}
