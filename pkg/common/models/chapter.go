package models

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model

	Name          string `gorm:"not null"`
	Description   string
	NumberOfPages int  `gorm:"not null"`
	Views         uint `gorm:"not null;default:0"`
	Volume        uint `gorm:"not null"`

	TitleID uint  `gorm:"not null"`
	Title   Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type ChapterOnModeration struct {
	gorm.Model

	Name          *string
	Description   *string
	NumberOfPages *int
	Volume        *uint

	TitleID *uint
	Title   *Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`

	TitleOnModerationID *uint
	TitleOnModeration   *TitleOnModeration `gorm:"foreignKey:TitleOnModerationID;references:id;constraint:OnDelete:CASCADE"`

	ExistingID *uint    `gorm:"unique"`
	Chapter    *Chapter `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (ChapterOnModeration) TableName() string {
	return "chapters_on_moderation"
}

func (c ChapterOnModeration) ToChapter() Chapter {
	chapter := Chapter{
		CreatorID:   c.CreatorID,
		TeamID:      c.TeamID,
		ModeratorID: c.ModeratorID,
	}

	if c.Name != nil {
		chapter.Name = *c.Name
	}
	if c.Description != nil {
		chapter.Description = *c.Description
	}
	if c.NumberOfPages != nil {
		chapter.NumberOfPages = *c.NumberOfPages
	}
	if c.Volume != nil {
		chapter.Volume = *c.Volume
	}
	if c.TitleID != nil {
		chapter.TitleID = *c.TitleID
	}

	return chapter
}
