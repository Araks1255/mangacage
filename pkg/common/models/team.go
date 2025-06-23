package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string `json:"description"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type TeamOnModeration struct {
	gorm.Model
	Name        sql.NullString `gorm:"unique"`
	Description string         `json:"description"`

	ExistingID *uint `gorm:"unique"`
	Team       *Team `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint  `gorm:"unique;not null"`
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:CASCADE"`
}

func (TeamOnModeration) TableName() string {
	return "teams_on_moderation"
}

type TeamDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        *string `json:"name"`
	Description *string `json:"description,omitempty"`
}

type TeamOnModerationDTO struct {
	TeamDTO
	Existing   string `json:"existing,omitempty"`
	ExistingID uint   `json:"existingId,omitempty"`
}
