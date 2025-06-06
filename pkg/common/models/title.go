package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Genres []Genre `gorm:"many2many:title_genres;constraint:OnDelete:CASCADE"`
}

type TitleOnModeration struct {
	gorm.Model
	Name        sql.NullString `gorm:"unique"`
	Description string

	ExistingID *uint `gorm:"unique"`
	Title      Title `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID *uint
	Author   *Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Genres []Genre `gorm:"many2many:title_on_moderation_genres;constraint:OnDelete:CASCADE"`
}

func (TitleOnModeration) TableName() string {
	return "titles_on_moderation"
}

type TitleDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Author   *string `json:"author,omitempty"`
	AuthorID *uint   `json:"authorId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Genres *pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`

	Views *int64 `json:"views,omitempty"`
}

type TitleOnModerationDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        *string `json:"name"`
	Description *string `json:"description,omitempty"`

	Author   *string `json:"author,omitempty"`
	AuthorID *uint   `json:"authorId,omitempty"`

	Genres *pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`
}
