package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Volume struct {
	gorm.Model
	Name        string
	Description string

	TitleID uint   `gorm:"not null"`
	Title   *Title `gorm:"foreignKey:TitleID;references:id"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type VolumeOnModeration struct {
	gorm.Model
	Name        sql.NullString
	Description string

	ExistingID *uint   `gorm:"unique"`
	Volume     *Volume `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	TitleID uint   `gorm:"not null"`
	Title   *Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:SET NULL"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (VolumeOnModeration) TableName() string {
	return "volumes_on_moderation"
}

type VolumeDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`
}

type VolumeOnModerationDTO struct {
	VolumeDTO
	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`
}
