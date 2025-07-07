package dto

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

type CreateChapterDTO struct { // Тут без тегов, ведь используется ручной парсинг
	ID                   *uint
	Name                 string
	Description          *string
	VolumeID             *uint
	VolumeOnModerationID *uint
	TeamID               uint
	Pages                [][]byte
}

type EditChapterDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ResponseChapterDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name          *string `json:"name"`
	Description   *string `json:"description,omitempty"`
	NumberOfPages *int    `json:"numberOfPages,omitempty"`

	Volume   *string `json:"volume,omitempty"`
	VolumeID *uint   `json:"volumeId,omitempty"`

	VolumeOnModeration   *string `json:"volumeOnModeration,omitempty"`
	VolumeOnModerationID *uint   `json:"volumeOnModerationId,omitempty"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	TitleOnModeration   *string `json:"titleOnModeration,omitempty"`
	TitleOnModerationID *uint   `json:"titleOnModerationID,omitempty"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Views *uint `json:"views,omitempty"`
}

func (c CreateChapterDTO) ToChapterOnModeration(userID uint) models.ChapterOnModeration {
	numberOfPages := len(c.Pages)
	return models.ChapterOnModeration{
		Name:                 &c.Name,
		Description:          c.Description,
		VolumeID:             c.VolumeID,
		VolumeOnModerationID: c.VolumeOnModerationID,
		TeamID:               c.TeamID,
		NumberOfPages:        &numberOfPages,
		CreatorID:            userID,
	}
}

func (c EditChapterDTO) ToChapterOnModeration(userID, existingID uint) models.ChapterOnModeration {
	return models.ChapterOnModeration{
		Name:        c.Name,
		Description: c.Description,
		CreatorID:   userID,
		ExistingID:  &existingID,
	}
}
