package models

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model

	Name          string
	Description   string
	NumberOfPages int
	Views         uint

	VolumeID uint    `gorm:"not null"`
	Volume   *Volume `gorm:"foreignKey:VolumeID;references:id" json:"-"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL" json:"-"`
}

type ChapterOnModeration struct {
	gorm.Model

	Name          *string
	Description   *string
	NumberOfPages *int

	ExistingID *uint    `gorm:"unique"`
	Chapter    *Chapter `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	VolumeID *uint
	Volume   *Volume `gorm:"foreignKey:VolumeID;references:id;constraint:OnDelete:SET NULL"`

	VolumeOnModerationID *uint
	VolumeOnModeration   *VolumeOnModeration `gorm:"foreignKey:VolumeOnModerationID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (ChapterOnModeration) TableName() string {
	return "chapters_on_moderation"
}
