package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName      string `gorm:"unique"`
	Password      string
	AboutYourself string
	TgUserID      int64

	TeamID sql.NullInt64
	Team   *Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Roles                    []Role    `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	TitlesUserIsSubscribedTo []Title   `gorm:"many2many:user_titles_subscribed_to;constraint:OnDelete:CASCADE"`
	ViewedChapters           []Chapter `gorm:"many2many:user_viewed_chapters;constraint:OnDelete:CASCADE"`

	FavoriteTitles   []Title   `gorm:"many2many:user_favorite_titles;constraint:OnDelete:CASCADE"`
	FavoriteChapters []Chapter `gorm:"many2many:user_favorite_chapters;constraint:OnDelete:CASCADE"`
	FavoriteGenres   []Genre   `gorm:"many2many:user_favorite_genres;constraint:OnDelete:CASCADE"`
}
