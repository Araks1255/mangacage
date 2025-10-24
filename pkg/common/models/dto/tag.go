package dto

import "github.com/Araks1255/mangacage/pkg/common/models"

type CreateTagDTO struct {
	Name string `json:"name" binding:"required,min=2,max=30"`
}

type ResponseTagDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (t CreateTagDTO) ToTagOnModeration(creatorID uint) models.TagOnModeration {
	return models.TagOnModeration{
		Name:      t.Name,
		CreatorID: creatorID,
	}
}
