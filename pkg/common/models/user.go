package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UserName      string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	AboutYourself *string
	TgUserID      *int64

	Visible     bool `gorm:"not null;default:false"`
	Verificated bool `gorm:"not null;default:false"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Roles                    []Role              `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	TitlesUserIsSubscribedTo []Title             `gorm:"many2many:user_titles_subscribed_to;constraint:OnDelete:CASCADE"`
	ViewedChapters           []UserViewedChapter `gorm:"foreignKey:UserID"`

	FavoriteTitles   []Title   `gorm:"many2many:user_favorite_titles;constraint:OnDelete:CASCADE"`
	FavoriteChapters []Chapter `gorm:"many2many:user_favorite_chapters;constraint:OnDelete:CASCADE"`
	FavoriteGenres   []Genre   `gorm:"many2many:user_favorite_genres;constraint:OnDelete:CASCADE"`
}

type UserOnModeration struct {
	gorm.Model

	UserName      *string
	AboutYourself *string

	ExistingID *uint `gorm:"unique"`
	Existing   *User `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`
}

func (UserOnModeration) TableName() string {
	return "users_on_moderation"
}
