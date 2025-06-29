package models

import (
	"mime/multipart"
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
	Visible       bool `gorm:"not null;default:false"`

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
	Password      *string
	AboutYourself *string
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

type UserOnModerationDTO struct {
	ID        uint      `json:"id" form:"-"`
	CreatedAt time.Time `json:"createdAt,omitempty" form:"-"`

	UserName      *string `json:"userName" form:"userName" binding:"required"`
	Password      *string `json:"password" form:"-" binding:"required"`
	AboutYourself *string `json:"aboutYourself" form:"aboutYourself"`

	Team   *string `json:"team,omitempty" form:"-"`
	TeamID *uint   `json:"teamId,omitempty" form:"-"`

	ProfilePicture *multipart.FileHeader `json:"-" form:"profilePicture"`
}

func (u UserOnModerationDTO) ToUserOnModeration(existingID *uint) UserOnModeration {
	return UserOnModeration{
		UserName:      u.UserName,
		Password:      u.Password,
		AboutYourself: u.AboutYourself,
		ExistingID:    existingID,
	}
}
