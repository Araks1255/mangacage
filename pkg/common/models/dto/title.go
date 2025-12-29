package dto

import (
	"mime/multipart"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"

	"github.com/lib/pq"
	"golang.org/x/text/unicode/norm"
)

type CreateTitleDTO struct {
	ID           *uint  `form:"id"` // Для upsert
	Name         string `form:"name" binding:"required,min=2,max=200"`
	EnglishName  string `form:"englishName" binding:"required,min=2,max=200"`
	OriginalName string `form:"originalName" binding:"required,min=2,max=200"`

	Description   *string `form:"description" binding:"max=2000"`
	AgeLimit      *int    `form:"ageLimit"`
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
	Name         *string `form:"name" binding:"omitempty,min=2,max=200"`
	EnglishName  *string `form:"englishName" binding:"omitempty,min=2,max=200"`
	OriginalName *string `form:"originalName" binding:"omitempty,min=2,max=200"`

	Description   *string `form:"description" binding:"omitempty,min=2,max=2000"`
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

	NumberOfRates *uint `json:"numberOfRates,omitempty"`
	SumOfRates    *uint `json:"sumOfRates,omitempty"`

	Author   *string `json:"author,omitempty"`
	AuthorID *uint   `json:"authorId,omitempty"`

	AuthorOnModeration   *string `json:"authorOnModeration,omitempty"`
	AuthorOnModerationID *uint   `json:"authorOnModerationId,omitempty"`

	Views *uint `json:"views,omitempty"`

	Genres   pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`
	Tags     pq.StringArray `json:"tags,omitempty" gorm:"type:TEXT[]"`
	Volumes  pq.Int64Array  `json:"volumes,omitempty" gorm:"type:BIGINT[]"`
	TeamsIDs pq.Int64Array  `json:"teamsIds,omitempty" gorm:"type:BIGINT[]"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`

	NumberOfViewedChapters *int64 `json:"numberOfViewedChapters,omitempty"`
	NumberOfChapters       *int64 `json:"numberOfChapters,omitempty"`

	UserRate       *int  `json:"userRate,omitempty"`
	FavoritedByMe  *bool `json:"favoritedByMe,omitempty"`
	MySubscription *bool `json:"mySubscription,omitempty"`
}

func (t CreateTitleDTO) ToTitleOnModeration(creatorID uint) *models.TitleOnModeration {
	formatedOriginalName := norm.NFC.String(t.OriginalName)

	var id uint
	if t.ID != nil {
		id = *t.ID
	}

	return &models.TitleOnModeration{
		Model:                gorm.Model{ID: id},
		CreatorID:            &creatorID,
		Name:                 &t.Name,
		EnglishName:          &t.EnglishName,
		OriginalName:         &formatedOriginalName,
		Description:          t.Description,
		AgeLimit:             t.AgeLimit,
		YearOfRelease:        &t.YearOfRelease,
		Type:                 &t.Type,
		TranslatingStatus:    &t.TranslatingStatus,
		PublishingStatus:     &t.PublishingStatus,
		AuthorID:             t.AuthorID,
		AuthorOnModerationID: t.AuthorOnModerationID,
	}
}

func (t EditTitleDTO) ToTitleOnModeration(creatorID, existingID uint) *models.TitleOnModeration {
	if t.OriginalName != nil {
		formatedOriginalName := norm.NFC.String(*t.OriginalName)
		t.OriginalName = &formatedOriginalName
	}

	return &models.TitleOnModeration{
		CreatorID:         &creatorID,
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
		ExistingID:        &existingID,
	}
}
