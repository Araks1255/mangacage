package moderation

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateGenreOnModeration(db *gorm.DB, userID uint) (uint, error) {
	genre := models.GenreOnModeration{
		Name: uuid.New().String(),
		CreatorID: userID,
	}

	if err := db.Create(&genre).Error; err != nil {
		return 0, err
	}

	return genre.ID, nil
}
