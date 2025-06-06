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

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Roles                    []Role              `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	TitlesUserIsSubscribedTo []Title             `gorm:"many2many:user_titles_subscribed_to;constraint:OnDelete:CASCADE"`
	ViewedChapters           []UserViewedChapter `gorm:"foreignKey:UserID"` // Это нужно по идее только для Preload. Я оставил для явности

	FavoriteTitles   []Title   `gorm:"many2many:user_favorite_titles;constraint:OnDelete:CASCADE"` // Тут как-нибудь вынесу всё в отдельные структуры и сделаю составной индекс уникальности на user_id и object_id  (чтобы на уровне бд нельзя было иметь две подписки на один объект)
	FavoriteChapters []Chapter `gorm:"many2many:user_favorite_chapters;constraint:OnDelete:CASCADE"`
	FavoriteGenres   []Genre   `gorm:"many2many:user_favorite_genres;constraint:OnDelete:CASCADE"`
}

type UserOnModeration struct {
	gorm.Model
	UserName      sql.NullString `gorm:"unique"`
	Password      string
	AboutYourself string
	TgUserID      int64

	ExistingID *uint `gorm:"unique"`
	User       *User `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`
}

func (UserOnModeration) TableName() string {
	return "users_on_moderation"
}

type UserDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	UserName      string  `json:"userName"`
	AboutYourself *string `json:"aboutYourself,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Roles *pq.StringArray `json:"roles,omitempty" gorm:"type:TEXT[]"`
}
