package moderation

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAuthorOnModeration(db *gorm.DB, userID uint) (uint, error) {
	author := models.AuthorOnModeration{
		Name:         uuid.New().String(),
		EnglishName:  uuid.New().String(),
		OriginalName: uuid.New().String(),
		CreatorID:    &userID,
	}

	if err := db.Create(&author).Error; err != nil {
		return 0, err
	}

	return author.ID, nil
}
