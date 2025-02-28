package models

import (
	"database/sql"

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
	Title   Title `gorm:"foreignKey:TitleID;references:id"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id"`
}
