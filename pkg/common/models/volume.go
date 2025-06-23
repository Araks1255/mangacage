package models

import (
	"time"

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
	Title   *Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:SET NULL"`

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

type VolumeDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`
}

type VolumeOnModerationDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`
}

func (v VolumeOnModerationDTO) ToVolumeOnModeration(creatorID uint, titleID, existingID *uint) VolumeOnModeration {
	return VolumeOnModeration{
		Name:        v.Name,
		Description: v.Description,
		TeamID:      v.TeamID,
		TitleID:     titleID,
		ExistingID:  existingID,
		CreatorID:   creatorID,
	}
}
