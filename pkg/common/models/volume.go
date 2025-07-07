package models

import (
	"gorm.io/gorm"
)

type Volume struct {
	gorm.Model

	Name        string `gorm:"not null"`
	Description string

	TitleID uint   `gorm:"not null"`
	Title   *Title `gorm:"foreignKey:TitleID;references:id"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint  `gorm:"not null"`
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type VolumeOnModeration struct {
	gorm.Model

	Name        *string
	Description *string

	ExistingID *uint   `gorm:"unique"`
	Volume     *Volume `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	TitleID *uint
	Title   *Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`

	TitleOnModerationID *uint
	TitleOnModeration   *TitleOnModeration `gorm:"foreignKey:TitleOnModerationID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (VolumeOnModeration) TableName() string {
	return "volumes_on_moderation"
}
