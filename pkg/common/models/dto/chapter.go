package dto

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

type CreateChapterDTO struct { // Тут без тегов, ведь используется ручной парсинг
	ID                  *uint
	Name                string
	Description         *string
	Volume              uint
	TitleID             *uint
	TitleOnModerationID *uint
	TeamID              uint
	Pages               [][]byte
}

type EditChapterDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Volume      *uint   `json:"volume"`
}

type ResponseChapterDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name          *string `json:"name"`
	Description   *string `json:"description,omitempty"`
	NumberOfPages *int    `json:"numberOfPages,omitempty"`
	Volume        uint    `json:"volume"`

	Title   *string `json:"title,omitempty"`
	TitleID *uint   `json:"titleId,omitempty"`

	TitleOnModeration   *string `json:"titleOnModeration,omitempty"`
	TitleOnModerationID *uint   `json:"titleOnModerationId,omitempty"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Views *uint `json:"views,omitempty"`
}

func (c CreateChapterDTO) ToChapterOnModeration(userID uint) models.ChapterOnModeration {
	numberOfPages := len(c.Pages)
	var id uint
	if c.ID != nil {
		id = *c.TitleOnModerationID
	}
	return models.ChapterOnModeration{
		Model: gorm.Model{
			ID: id,
		},
		Name:                &c.Name,
		Description:         c.Description,
		Volume:              &c.Volume,
		TitleID:             c.TitleID,
		TitleOnModerationID: c.TitleOnModerationID,
		TeamID:              c.TeamID,
		NumberOfPages:       &numberOfPages,
		CreatorID:           userID,
	}
}

func (c EditChapterDTO) ToChapterOnModeration(userID, existingID uint) models.ChapterOnModeration {
	return models.ChapterOnModeration{
		Name:        c.Name,
		Description: c.Description,
		Volume:      c.Volume,
		CreatorID:   userID,
		ExistingID:  &existingID,
	}
}
