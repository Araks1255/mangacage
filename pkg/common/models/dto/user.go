package dto

import (
	"mime/multipart"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/lib/pq"
)

type CreateUserDTO struct {
	UserName      string  `json:"userName" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	AboutYourself *string `json:"aboutYourself"`
}

type EditUserDTO struct {
	UserName       *string               `form:"userName"`
	AboutYourself  *string               `form:"aboutYourself"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
}

type ResponseUserDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	UserName      string  `json:"userName"`
	AboutYourself *string `json:"aboutYourself,omitempty"`
	Visible       bool

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`

	Existing   *string `json:"existing,omitempty"`
	ExistingID *uint   `json:"existingId,omitempty"`

	Roles *pq.StringArray `json:"roles,omitempty" gorm:"type:TEXT[]"`
}

func (u CreateUserDTO) ToUser() models.User {
	return models.User{
		UserName:      u.UserName,
		Password:      u.Password,
		AboutYourself: u.AboutYourself,
		Verificated:   false,
	}
}

func (u EditUserDTO) ToUserOnModeration(existingID uint) models.UserOnModeration {
	return models.UserOnModeration{
		UserName:      u.UserName,
		AboutYourself: u.AboutYourself,
		ExistingID:    &existingID,
	}
}
