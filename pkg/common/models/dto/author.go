package dto

import (
	"github.com/Araks1255/mangacage/pkg/common/models"

	"golang.org/x/text/unicode/norm"
)

type CreateAuthorDTO struct {
	Name         string  `json:"name" binding:"required,min=2,max=50"`
	EnglishName  string  `json:"englishName" binding:"required,min=2,max=50"`
	OriginalName string  `json:"originalName" binding:"required,min=2,max=50"`
	About        *string `json:"about" binding:"omitempty,max=400"`
}

type ResponseAuthorDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	EnglishName  string  `json:"englishName,omitempty"`
	OriginalName string  `json:"originalName,omitempty"`
	About        *string `json:"about,omitempty"`
}

func (a CreateAuthorDTO) ToAuthorOnModeration(creatorID uint) models.AuthorOnModeration {
	formatedOriginalName := norm.NFC.String(a.OriginalName)
	return models.AuthorOnModeration{
		Name:         a.Name,
		EnglishName:  a.EnglishName,
		OriginalName: formatedOriginalName,
		About:        a.About,
		CreatorID:    &creatorID,
	}
}
