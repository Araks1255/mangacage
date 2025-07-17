package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	Name        string `gorm:"not null"`
	Description string

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type TeamOnModeration struct {
	gorm.Model
	Name        *string
	Description *string

	ExistingID *uint `gorm:"unique"`
	Team       *Team `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint  `gorm:"unique;not null"`
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (TeamOnModeration) TableName() string {
	return "teams_on_moderation"
}
