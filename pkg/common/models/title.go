package models

import (
	"mime/multipart"
	"time"

	"github.com/lib/pq"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model

	Name         string `gorm:"unique"`
	EnglishName  string `gorm:"unique"`
	OriginalName string `gorm:"unique"`

	Description   *string
	AgeLimit      int    `gorm:"not null"`
	YearOfRelease int    `gorm:"not null"`
	Type          string `gorm:"type:title_type"`

	TranslatingStatus string `gorm:"type:title_translating_status;default:'free'"`
	PublishingStatus  string `gorm:"type:title_publishing_status;default:'unknown'"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID uint
	Author   Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Views uint `gorm:"default:0;not null"`

	Genres []Genre `gorm:"many2many:title_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_tags;constraint:OnDelete:CASCADE"`
}

type TitleOnModeration struct {
	gorm.Model

	Name         *string `gorm:"unique"`
	EnglishName  *string `gorm:"unique"`
	OriginalName *string `gorm:"unique"`

	Description   *string
	AgeLimit      *int
	YearOfRelease *int
	Type          *string `gorm:"type:title_type"`

	TranslatingStatus *string `gorm:"type:title_translating_status"`
	PublishingStatus  *string `gorm:"type:title_publishing_status"`

	ExistingID *uint `gorm:"unique"`
	Title      Title `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	AuthorID *uint
	Author   *Author `gorm:"foreignKey:AuthorID;references:id;constraint:OnDelete:SET NULL"`

	TeamID *uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:SET NULL"`

	Genres []Genre `gorm:"many2many:title_on_moderation_genres;constraint:OnDelete:CASCADE"`
	Tags   []Tag   `gorm:"many2many:title_on_moderation_tags;constraint:OnDelete:CASCADE"`
}

func (TitleOnModeration) TableName() string {
	return "titles_on_moderation"
}

type TitleDTO struct {
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

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Views *uint `json:"views,omitempty"`

	Genres *pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`
	Tags   *pq.StringArray `json:"tags,omitempty" gorm:"type:TEXT[]"`
}

type TitleOnModerationDTO struct {
	ID        uint       `json:"id" form:"-"`
	CreatedAt *time.Time `json:"createdAt,omitempty" form:"-"`

	Name         *string `json:"name" form:"name" binding:"required"`
	EnglishName  *string `json:"englishName,omitempty" form:"englishName" binding:"required"`
	OriginalName *string `json:"originalName,omitempty" form:"originalName" binding:"required"`

	Description   *string `json:"description,omitempty" form:"description"`
	AgeLimit      *int    `json:"ageLimit,omitempty" form:"ageLimit" binding:"required"`
	YearOfRelease *int    `json:"yearOfRelease,omitempty" form:"yearOfRelease" binding:"required"`
	Type          *string `json:"type,omitempty" form:"type" binding:"required"`

	TranslatingStatus *string `json:"translatingStatus,omitempty" form:"translatingStatus"`
	PublishingStatus  *string `json:"publishingStatus,omitempty" form:"publishingStatus" binding:"required"`

	Author   *string `json:"author,omitempty" form:"-"`
	AuthorID *uint   `json:"authorId,omitempty" form:"authorId" binding:"required"`

	GenresIDs []uint          `json:"-" form:"genresIds" binding:"required"`
	Genres    *pq.StringArray `json:"genres,omitempty" form:"-" gorm:"type:TEXT[]"`

	TagsIDs []uint `json:"-" form:"tagsIds" binding:"required"`

	Cover *multipart.FileHeader `json:"-" form:"cover" binding:"required"` // Ограничение выставить

	Existing   *string `json:"existing,omitempty" form:"-"`
	ExistingID *uint   `json:"existingId,omitempty" form:"-"`
}

func (t TitleOnModerationDTO) ToTitleOnModeration(creatorID uint, existingID *uint) TitleOnModeration {
	var res TitleOnModeration

	if t.OriginalName != nil {
		formatedOriginalName := norm.NFC.String(*t.OriginalName)
		res.OriginalName = &formatedOriginalName
	}

	res.ID = t.ID
	res.CreatorID = creatorID
	res.ExistingID = existingID

	res.Name = t.Name
	res.EnglishName = t.EnglishName

	res.Description = t.Description
	res.AgeLimit = t.AgeLimit
	res.YearOfRelease = t.YearOfRelease
	res.Type = t.Type

	res.TranslatingStatus = t.TranslatingStatus
	res.PublishingStatus = t.PublishingStatus

	res.AuthorID = t.AuthorID

	return res
}
