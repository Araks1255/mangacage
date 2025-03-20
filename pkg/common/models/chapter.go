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

	VolumeID uint
	Volume   Volume `gorm:"foreignKey:VolumeID;references:id;OnDelete:SET NULL" json:"-"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL" json:"-"`
}
