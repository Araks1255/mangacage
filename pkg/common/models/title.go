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
	AgeLimit      int    `gorm:"not null"`
	YearOfRelease int    `gorm:"not null"`
	Type          string `gorm:"type:title_type;not null"`

	TranslatingStatus string `gorm:"type:title_translating_status;default:'free'"`
	PublishingStatus  string `gorm:"type:title_publishing_status;default:'unknown'"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	Views uint `gorm:"default:0;not null"`

	SumOfRates    uint `gorm:"not null"`
	NumberOfRates uint `gorm:"not null"`

	Genres []Genre `gorm:"many2many:title_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_tags;constraint:OnDelete:CASCADE"`
	Teams  []Team  `gorm:"many2many:title_teams;constraint:OnDelete:CASCADE"`
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

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID *uint
	Author   *Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorOnModerationID *uint
	AuthorOnModeration   *AuthorOnModeration `gorm:"foreignKey:AuthorOnModerationID;references:id;constraint:OnDelete:CASCADE"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Genres []Genre `gorm:"many2many:title_on_moderation_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_on_moderation_tags;constraint:OnDelete:CASCADE"`
}

func (TitleOnModeration) TableName() string {
	return "titles_on_moderation"
}
