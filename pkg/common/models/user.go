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
	Team   Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Roles                    []Role    `gorm:"many2many:user_roles;OnDelete:SET NULL"`
	TitlesUserIsSubscribedTo []Title   `gorm:"many2many:user_titles_subscribed_to;OnDelete:SET NULL"`
	ViewedChapters           []Chapter `gorm:"many2many:user_viewed_chapters;OnDelete:SET NULL"`

	FavoriteTitles   []Title   `gorm:"many2many:user_favorite_titles;OnDelete:SET NULL"`
	FavoriteChapters []Chapter `gorm:"many2many:user_favorite_chapters;OnDelete:SET NULL"`
	FavoriteGenres   []Genre   `gorm:"many2many:user_favorite_genres;OnDelete:SET NULL"`
}
