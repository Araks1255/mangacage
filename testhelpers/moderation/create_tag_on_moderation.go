package moderation

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTagOnModeration(db *gorm.DB, userID uint) (uint, error) {
	tag := models.TagOnModeration{
		Name:      uuid.New().String(),
		CreatorID: userID,
	}

	if err := db.Create(&tag).Error; err != nil {
		return 0, err
	}

	return tag.ID, nil
}
