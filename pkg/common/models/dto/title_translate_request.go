package dto

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

type CreateTitleTranslateRequestDTO struct {
	Message *string `json:"message" binding:"omitempty,max=100"`
	TeamID  uint    `json:"teamId" binding:"required"`
	TitleID uint    `json:"titleId" binding:"required"`
}

type ResponseTitleTranslateRequestDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Message *string `json:"message,omitempty"`

	Team   string `json:"team"`
	TeamID uint   `json:"teamId"`

	Title   string `json:"title"`
	TitleID uint   `json:"titleId"`
}

func (ttr CreateTitleTranslateRequestDTO) ToTitleTranslateRequest() models.TitleTranslateRequest {
	return models.TitleTranslateRequest{
		Message: ttr.Message,
		TeamID:  ttr.TeamID,
		TitleID: ttr.TitleID,
	}
}
