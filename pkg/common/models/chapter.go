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

	VolumeID uint
	Volume   *Volume `gorm:"foreignKey:VolumeID;references:id" json:"-"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL" json:"-"`
}

type ChapterOnModeration struct {
	gorm.Model
	Name          string
	Description   string
	NumberOfPages int

	ExistingID sql.NullInt64 `gorm:"unique"`
	Chapter    *Chapter      `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	VolumeID sql.NullInt64
	Volume   *Volume `gorm:"foreignKey:VolumeID;references:id;OnDelete:SET NULL" json:"-"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL" json:"-"`
}

func (ChapterOnModeration) TableName() string {
	return "chapters_on_moderation"
}

type ChapterDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	NumberOfPages int    `json:"numberOfPages,omitempty"`

	Volume   string `json:"volume,omitempty"`
	VolumeID uint   `json:"volumeId,omitempty"`

	Title   string `json:"title,omitempty"`
	TitleID uint   `json:"titleId,omitempty"`
}

type ChapterOnModerationDTO struct {
	ChapterDTO
	Existing   string `json:"existing,omitempty"`
	ExistingID uint   `json:"existingId,omitempty"`
}
