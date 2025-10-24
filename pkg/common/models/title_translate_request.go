package models

import "gorm.io/gorm"

type TitleTranslateRequest struct {
	gorm.Model

	Message *string

	TeamID uint `gorm:"not null"`
	Team   Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`

	TitleID uint  `gorm:"not null"`
	Title   Title `gorm:"foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint
}

func (TitleTranslateRequest) TableName() string {
	return "titles_translate_requests"
}
