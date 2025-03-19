package models

import (
	"database/sql"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TitleOnModeration struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string

	ExistingID sql.NullInt64
	Title      Title `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id;OnDelete:SET NULL"`

	TeamID sql.NullInt64
	Team   Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Genres pq.StringArray `gorm:"type:text[]"`
}

func (TitleOnModeration) TableName() string {
	return "titles_on_moderation"
}

type VolumeOnModeration struct {
	gorm.Model
	Name        string
	Description string

	ExistingID sql.NullInt64
	Volume     Volume `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	TitleID uint
	Title   Title `gorm:"foreignKey:TitleID;references:id;OnDelete:SET NULL"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`
}

func (VolumeOnModeration) TableName() string {
	return "volumes_on_moderation"
}

type ChapterOnModeration struct {
	gorm.Model
	Name          string
	Description   string
	NumberOfPages int

	ExistingID sql.NullInt64
	Chapter    Chapter `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	VolumeID uint
	Volume   Volume `gorm:"foreignKey:VolumeID;references:id;OnDelete:SET NULL" json:"-"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL" json:"-"`
}

func (ChapterOnModeration) TableName() string {
	return "chapters_on_moderation"
}

type UserOnModeration struct {
	gorm.Model
	UserName      string `gorm:"unique"`
	Password      string
	AboutYourself string
	TgUserID      int64

	ExistingID sql.NullInt64
	User       User `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	TeamID sql.NullInt64
	Team   Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Roles pq.StringArray `gorm:"type:text[]"`
}

func (UserOnModeration) TableName() string {
	return "users_on_moderation"
}

type TeamOnModeration struct {
	gorm.Model
	Name        string `json:"name" binding:"required" gorm:"unique"`
	Description string `json:"description"`

	ExistingID sql.NullInt64
	Team       Team `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`
}

func (TeamOnModeration) TableName() string {
	return "teams_on_moderation"
}
