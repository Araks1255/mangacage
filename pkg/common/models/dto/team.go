package dto

import (
	"mime/multipart"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

type CreateTeamDTO struct {
	ID          *uint                 `form:"id"`
	Name        string                `form:"name" binding:"required"`
	Description *string               `form:"description" binding:"required"`
	Cover       *multipart.FileHeader `form:"cover" binding:"required"`
}

type EditTeamDTO struct {
	Name        *string               `form:"name"`
	Description *string               `form:"description"`
	Cover       *multipart.FileHeader `form:"cover"`
}

type ResponseTeamDTO struct {
	ID          uint       `json:"id"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	Name        *string    `json:"name"`
	Description *string    `json:"description,omitempty"`
	Existing    *string    `json:"existing,omitempty"`
	ExistingID  *uint      `json:"existingId,omitempty"`
}

func (t CreateTeamDTO) ToTeamOnModeration(userID uint) models.TeamOnModeration {
	var id uint
	if t.ID != nil {
		id = *t.ID
	}
	return models.TeamOnModeration{
		Model: gorm.Model{
			ID: id,
		},
		Name:        &t.Name,
		Description: t.Description,
		CreatorID:   userID,
	}
}

func (t EditTeamDTO) ToTeamOnModeration(userID, existingID uint) models.TeamOnModeration {
	return models.TeamOnModeration{
		Name:        t.Name,
		Description: t.Description,
		CreatorID:   userID,
		ExistingID:  &existingID,
	}
}
