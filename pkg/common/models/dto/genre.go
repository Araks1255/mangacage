package dto

import "github.com/Araks1255/mangacage/pkg/common/models"

type CreateGenreDTO struct {
	Name string `json:"name" binding:"required,min=2,max=30"`
}

type ResponseGenreDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (g CreateGenreDTO) ToGenreOnModeration(creatorID uint) models.GenreOnModeration {
	return models.GenreOnModeration{
		Name:      g.Name,
		CreatorID: creatorID,
	}
}
