package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID           uint `gorm:"primaryKey;autoIncrement:true"`
	Name         string
	EnglishName  string
	OriginalName string `gorm:"unique"`
	About        string
}

type AuthorDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	EnglishName  string  `json:"englishName"`
	OriginalName string  `json:"originalName"`
	About        *string `json:"about,omitempty"`
}

type AuthorOnModeration struct {
	gorm.Model

	Name         string
	EnglishName  string
	OriginalName string `gorm:"unique"`
	About        *string

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:CASCADE"`
}

func (AuthorOnModeration) TableName() string {
	return "authors_on_moderation"
}

type AuthorOnModerationDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Name         string  `json:"name" binding:"required"`
	EnglishName  string  `json:"englishName" binding:"required"`
	OriginalName string  `json:"originalName" binding:"required"`
	About        *string `json:"about,omitempty"`

	CreatorID *uint `json:"creatorId,omitempty"`
}

func (a AuthorOnModerationDTO) ToAuthorOnModeration(creatorID uint) AuthorOnModeration {
	return AuthorOnModeration{
		Name:         a.Name,
		EnglishName:  a.EnglishName,
		OriginalName: a.OriginalName,
		About:        a.About,
		CreatorID:    creatorID,
	}
}
