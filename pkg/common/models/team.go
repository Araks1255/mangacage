package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `json:"name" binding:"required" gorm:"unique"`
	Description string `json:"description"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;OnDelete:SET NULL"`

	ModeratorID sql.NullInt64
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;OnDelete:SET NULL"`
}

type TeamOnModeration struct {
	gorm.Model
	Name        sql.NullString `json:"name" binding:"required" gorm:"unique"`
	Description string         `json:"description"`

	ExistingID sql.NullInt64 `gorm:"unique"`
	Team       *Team         `gorm:"foreignKey:ExistingID;references:id;OnDelete:CASCADE"`

	CreatorID uint  `gorm:"unique"`
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;OnDelete:CASCADE"`
}

func (TeamOnModeration) TableName() string {
	return "teams_on_moderation"
}

type TeamDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Leader   string `json:"leader,omitempty"`
	LeaderID uint   `json:"leaderId,omitempty"`
}
