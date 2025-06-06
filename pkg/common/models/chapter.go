package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	Name          string
	Description   string
	NumberOfPages int

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
	Name          sql.NullString
	Description   string
	NumberOfPages int

	ExistingID *uint    `gorm:"unique"`
	Chapter    *Chapter `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	VolumeID uint    `gorm:"not null"`
	Volume   *Volume `gorm:"foreignKey:VolumeID;references:id;constraint:OnDelete:SET NULL"`

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

type ChapterDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	NumberOfPages *int    `json:"numberOfPages,omitempty"`

	Volume   *string `json:"volume,omitempty"`
	VolumeID *uint   `json:"volumeId,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`
}

type ChapterOnModerationDTO struct {
	ChapterDTO
	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`
}
