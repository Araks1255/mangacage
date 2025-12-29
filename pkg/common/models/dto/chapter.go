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
	NumberOfPages       int

	WebtoonMode        bool
	DisableCompression bool
	NeedsCompression   bool

	PagesSize       int64
	PagesResolution int64
}

type EditChapterDTO struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=35"`
	Description *string `json:"description" binding:"omitempty,max=100"`
	Volume      *uint   `json:"volume"`
}

type ResponseChapterDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name          *string `json:"name"`
	Description   *string `json:"description,omitempty"`
	NumberOfPages *int    `json:"numberOfPages,omitempty"`
	Volume        *uint   `json:"volume,omitempty"`
	WebtoonMode   bool    `json:"webtoonMode"`

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
	var id uint
	if c.ID != nil && *c.ID != 0 {
		id = *c.TitleOnModerationID
	}
	if c.TitleID != nil && *c.TitleID == 0 {
		c.TitleID = nil
	}
	if c.TitleOnModerationID != nil && *c.TitleOnModerationID == 0 {
		c.TitleOnModerationID = nil
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
		CreatorID:           &userID,
		WebtoonMode:         c.WebtoonMode,
	}
}

func (c EditChapterDTO) ToChapterOnModeration(userID, existingID uint) models.ChapterOnModeration {
	return models.ChapterOnModeration{
		Name:        c.Name,
		Description: c.Description,
		Volume:      c.Volume,
		CreatorID:   &userID,
		ExistingID:  &existingID,
	}
}
