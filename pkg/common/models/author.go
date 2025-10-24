package models

import (
	"gorm.io/gorm"
)

type Author struct {
	ID           uint   `gorm:"primaryKey;autoIncrement:true"`
	Name         string `gorm:"not null"`
	EnglishName  string `gorm:"not null"`
	OriginalName string `gorm:"unique;not null"`
	About        *string

	ModeratorID *uint
}

type AuthorOnModeration struct {
	gorm.Model

	Name         string `gorm:"not null"`
	EnglishName  string `gorm:"not null"`
	OriginalName string `gorm:"unique;not null"`
	About        *string

	CreatorID *uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
}

func (AuthorOnModeration) TableName() string {
	return "authors_on_moderation"
}

func (a AuthorOnModeration) ToAuthor() *Author {
	return &Author{
		Name:         a.Name,
		EnglishName:  a.EnglishName,
		OriginalName: a.OriginalName,
		About:        a.About,
		ModeratorID:  a.ModeratorID,
	}
}
