package models

import (
	"database/sql"
	"mime/multipart"
	"time"

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
	Views         *uint   `json:"views,omitempty"`

	Volume   *string `json:"volume,omitempty"`
	VolumeID *uint   `json:"volumeId,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`
}

type ChapterOnModerationDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name          string  `json:"name" form:"name" binding:"required"`
	Description   *string `json:"description,omitempty" form:"name"`
	NumberOfPages *int    `json:"numberOfPages,omitempty" form:"-"`

	Volume   *string `json:"volume,omitempty" form:"-"`
	VolumeID *uint   `json:"volumeId,omitempty" form:"volumeId"`

	Title   *string `json:"title,omitempty" form:"-"`
	TitleID *uint   `json:"titleId,omitempty" form:"-"`

	Existing   *string `json:"existing,omitempty" form:"-"`
	ExistingID *uint   `json:"existingId,omitempty" form:"existingId"`

	Pages []*multipart.FileHeader `json:"-" form:"pages" binding:"required" gorm:"-"`
}
