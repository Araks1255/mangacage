package testhelpers

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func CreateTitleTranslateRequest(db *gorm.DB, titleID, teamID uint, message ...string) (uint, error) {
	request := models.TitleTranslateRequest{
		TitleID: titleID,
		TeamID:  teamID,
	}

	if len(message) != 0 {
		request.Message = &message[0]
	}

	if err := db.Create(&request).Error; err != nil {
		return 0, err
	}

	return request.ID, nil
}
