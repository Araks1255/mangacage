package testhelpers

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAuthor(db *gorm.DB) (uint, error) {
	author := models.Author{
		Name:         uuid.New().String(),
		EnglishName:  uuid.New().String(),
		OriginalName: uuid.New().String(),
	}

	if err := db.Create(&author).Error; err != nil {
		return 0, err
	}

	return author.ID, nil
}
