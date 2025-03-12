package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model

	Name          string
	Description   string
	NumberOfPages int
	OnModeration  bool

	VolumeID uint
	Volume   Volume `gorm:"foreignKey:VolumeID;references:id;OnDelete:SET NULL" json:"-"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL" json:"-"`
}
