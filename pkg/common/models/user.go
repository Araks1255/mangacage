package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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

type UserDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	UserName      string `json:"userName"`
	AboutYourself string `json:"aboutYourself,omitempty"`

	Team   string `json:"team,omitempty"`
	TeamID uint   `json:"teamId,omitempty"`

	Roles pq.StringArray `json:"roles,omitempty" gorm:"type:TEXT[]"`
}

type UserOnModeration struct {
	gorm.Model
	UserName      sql.NullString `gorm:"unique"`
	Password      string
	AboutYourself string
	TgUserID      int64

	ExistingID sql.NullInt64 `gorm:"unique"`
	User       *User         `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	TeamID sql.NullInt64
	Team   *Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Roles pq.StringArray `gorm:"type:text[]"`
}

func (UserOnModeration) TableName() string {
	return "users_on_moderation"
}
