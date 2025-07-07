package dto

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

type CreateVolumeDTO struct {
	ID          *uint   `json:"id"` // Для upsert
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`

	TitleOnModerationID *uint `json:"titleOnModerationId" binding:"required_without=TitleID"`
	TitleID             *uint `json:"titleId" binding:"required_without=TitleOnModerationID"`

	TeamID *uint
}

type EditVolumeDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ResponseVolumeDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name        *string `json:"name"`
	Description *string `json:"description,omitempty"`

	TitleOnModeration   *string `json:"titleOnModeration,omitempty"`
	TitleOnModerationID *uint   `json:"titleOnModerationId,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`
}

func (v CreateVolumeDTO) ToVolumeOnModeration(userID uint) models.VolumeOnModeration {
	return models.VolumeOnModeration{
		Name:                &v.Name,
		Description:         v.Description,
		TitleID:             v.TitleID,
		TitleOnModerationID: v.TitleOnModerationID,
		TeamID:              v.TeamID,
		CreatorID:           userID,
	}
}

func (v EditVolumeDTO) ToVolumeOnModeration(userID, existingID uint) models.VolumeOnModeration {
	return models.VolumeOnModeration{
		Name:        v.Name,
		Description: v.Description,
		ExistingID:  &existingID,
	}
}
