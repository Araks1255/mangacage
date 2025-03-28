package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `json:"name" binding:"required" gorm:"unique"`
	Description string `json:"description"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`
}
