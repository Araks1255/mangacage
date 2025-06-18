package models

import (
	"time"

	"gorm.io/gorm"
)

type TitleTranslateRequest struct {
	gorm.Model

	Message *string

	TeamID uint
	Team   Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`

	TitleID uint
	Title   Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`
}

type TitleTranslateRequestDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Message *string `json:"message"`

	Team   string `json:"team"`
	TeamID uint   `json:"teamId"`

	Title   string `json:"title"`
	TitleID uint   `json:"titleId"`
}
