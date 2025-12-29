package models

import (
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model

	Name         string `gorm:"not null"`
	EnglishName  string `gorm:"not null"`
	OriginalName string `gorm:"not null"`

	Description   *string
	AgeLimit      int
	YearOfRelease int    `gorm:"not null"`
	Type          string `gorm:"type:title_type;not null"`

	TranslatingStatus string `gorm:"type:title_translating_status;default:'free'"`
	PublishingStatus  string `gorm:"type:title_publishing_status;default:'unknown'"`

	CreatorID *uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	EditorID *uint
	Editor   *User `gorm:"foreignKey:EditorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID uint   `gorm:"not null"`
	Author   Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	Views uint `gorm:"default:0;not null"`

	SumOfRates    uint `gorm:"not null"`
	NumberOfRates uint `gorm:"not null"`

	NumberOfChapters uint `gorm:"not null"`

	Genres []Genre `gorm:"many2many:title_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_tags;constraint:OnDelete:CASCADE"`
	Teams  []Team  `gorm:"many2many:title_teams;constraint:OnDelete:CASCADE"`

	ModeratorID *uint

	Hidden bool `gorm:"not null;default:false"`

	CoverPath string `gorm:"not null"`
}

type TitleOnModeration struct {
	gorm.Model

	Name         *string
	EnglishName  *string
	OriginalName *string

	Description   *string
	AgeLimit      *int
	YearOfRelease *int
	Type          *string `gorm:"type:title_type"`

	TranslatingStatus *string `gorm:"type:title_translating_status"`
	PublishingStatus  *string `gorm:"type:title_publishing_status"`

	ExistingID *uint `gorm:"unique"`
	Title      Title `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID *uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID *uint
	Author   *Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorOnModerationID *uint
	AuthorOnModeration   *AuthorOnModeration `gorm:"foreignKey:AuthorOnModerationID;references:id;constraint:OnDelete:CASCADE"`

	Genres []Genre `gorm:"many2many:title_on_moderation_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_on_moderation_tags;constraint:OnDelete:CASCADE"`

	ModeratorID *uint

	CoverPath *string
}

func (TitleOnModeration) TableName() string {
	return "titles_on_moderation"
}

func (t TitleOnModeration) ToTitle() *Title {
	title := &Title{
		ModeratorID: t.ModeratorID,
		Genres:      t.Genres,
		Tags:        t.Tags,
	}

	if t.Name != nil {
		title.Name = *t.Name
	}
	if t.EnglishName != nil {
		title.EnglishName = *t.EnglishName
	}
	if t.OriginalName != nil {
		title.OriginalName = *t.OriginalName
	}
	if t.Description != nil {
		title.Description = t.Description
	}
	if t.AgeLimit != nil {
		title.AgeLimit = *t.AgeLimit
	}
	if t.YearOfRelease != nil {
		title.YearOfRelease = *t.YearOfRelease
	}
	if t.Type != nil {
		title.Type = *t.Type
	}
	if t.TranslatingStatus != nil {
		title.TranslatingStatus = *t.TranslatingStatus
	}
	if t.PublishingStatus != nil {
		title.PublishingStatus = *t.PublishingStatus
	}
	if t.AuthorID != nil {
		title.AuthorID = *t.AuthorID
	}

	if t.ExistingID == nil {
		title.CreatorID = t.CreatorID
	} else {
		title.EditorID = t.CreatorID
	}

	return title
}

func (t *TitleOnModeration) SetID(id uint) {
	t.ID = id
}
