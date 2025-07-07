package dto

import (
	"mime/multipart"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"github.com/lib/pq"
	"golang.org/x/text/unicode/norm"
)

type CreateTitleDTO struct {
	ID           *uint  `form:"id"` // Для upsert
	Name         string `form:"name" binding:"required"`
	EnglishName  string `form:"englishName" binding:"required"`
	OriginalName string `form:"originalName" binding:"required"`

	Description   *string `form:"description"`
	AgeLimit      int     `form:"ageLimit" binding:"required"`
	YearOfRelease int     `form:"yearOfRelease" binding:"required"`
	Type          string  `form:"type" binding:"required"`

	TranslatingStatus string `form:"translatingStatus,default=free"`
	PublishingStatus  string `form:"publishingStatus" binding:"required"`

	AuthorID             *uint `form:"authorId" binding:"required_without=AuthorOnModerationID"`
	AuthorOnModerationID *uint `form:"authorOnModerationId" binding:"required_without=AuthorID"`

	GenresIDs []uint `form:"genresIds" binding:"required"`

	TagsIDs []uint `form:"tagsIds" binding:"required"`

	Cover *multipart.FileHeader `form:"cover" binding:"required"`
}

type EditTitleDTO struct {
	Name         *string `form:"name"`
	EnglishName  *string `form:"englishName"`
	OriginalName *string `form:"originalName"`

	Description   *string `form:"description"`
	AgeLimit      *int    `form:"ageLimit"`
	YearOfRelease *int    `form:"yearOfRelease"`
	Type          *string `form:"type"`

	TranslatingStatus *string `form:"translatingStatus"`
	PublishingStatus  *string `form:"publishingStatus"`

	AuthorID *uint `form:"authorId"`

	GenresIDs []uint `form:"genresIds"`
	TagsIDs   []uint `form:"tagsIds"`

	ExistingID uint

	Cover *multipart.FileHeader `form:"cover"`
}

type ResponseTitleDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	Name         *string `json:"name"`
	EnglishName  *string `json:"englishName,omitempty"`
	OriginalName *string `json:"originalName,omitempty"`

	Description   *string `json:"description,omitempty"`
	AgeLimit      *int    `json:"ageLimit,omitempty"`
	YearOfRelease *int    `json:"yearOfRelease,omitempty"`
	Type          *string `json:"type,omitempty"`

	TranslatingStatus *string `json:"translatingStatus,omitempty"`
	PublishingStatus  *string `json:"publishingStatus,omitempty"`

	Author   *string `json:"author,omitempty"`
	AuthorID *uint   `json:"authorId,omitempty"`

	AuthorOnModeration   *string `json:"authorOnModeration,omitempty"`
	AuthorOnModerationID *uint   `json:"authorOnModerationId,omitempty"`

	Views *uint    `json:"views,omitempty"`
	Rate  *float64 `json:"rate,omitempty"`

	Genres pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`
	Tags   pq.StringArray `json:"tags,omitempty" gorm:"type:TEXT[]"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`

	QuantityOfViewedChapters *int64 `json:"quantityOfViewedChapters,omitempty"`
	UserRate                 *int   `json:"userRate,omitempty"`
	QuantityOfChapters       *int64 `json:"quantityOfChapters,omitempty"`
	CanEdit                  *bool  `json:"canEdit,omitempty"`
}

func (t CreateTitleDTO) ToTitleOnModeration(creatorID uint) models.TitleOnModeration {
	formatedOriginalName := norm.NFC.String(t.OriginalName)
	return models.TitleOnModeration{
		CreatorID:            creatorID,
		Name:                 &t.Name,
		EnglishName:          &t.EnglishName,
		OriginalName:         &formatedOriginalName,
		Description:          t.Description,
		AgeLimit:             &t.AgeLimit,
		YearOfRelease:        &t.YearOfRelease,
		Type:                 &t.Type,
		TranslatingStatus:    &t.TranslatingStatus,
		PublishingStatus:     &t.PublishingStatus,
		AuthorID:             t.AuthorID,
		AuthorOnModerationID: t.AuthorOnModerationID,
	}
}

func (t EditTitleDTO) ToTitleOnModeration(creatorID, existingID uint) models.TitleOnModeration {
	if t.OriginalName != nil {
		formatedOriginalName := norm.NFC.String(*t.OriginalName)
		t.OriginalName = &formatedOriginalName
	}
	return models.TitleOnModeration{
		CreatorID:         creatorID,
		Name:              t.Name,
		EnglishName:       t.EnglishName,
		OriginalName:      t.OriginalName,
		Description:       t.Description,
		AgeLimit:          t.AgeLimit,
		YearOfRelease:     t.YearOfRelease,
		Type:              t.Type,
		TranslatingStatus: t.TranslatingStatus,
		PublishingStatus:  t.PublishingStatus,
		AuthorID:          t.AuthorID,
	}
}
