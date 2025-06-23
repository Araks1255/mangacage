package moderation

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateVolumeOnModerationOptions struct {
	ExistingID uint
}

func CreateVolumeOnModeration(db *gorm.DB, titleID, teamID, userID uint, opts ...CreateVolumeOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	name := uuid.New().String()

	volume := models.VolumeOnModeration{
		Name:      &name,
		TitleID:   &titleID,
		CreatorID: userID,
		TeamID:    &teamID,
	}

	if len(opts) != 0 && opts[0].ExistingID != 0 {
		volume.ExistingID = &opts[0].ExistingID
	}

	if result := db.Create(&volume); result.Error != nil {
		return 0, result.Error
	}

	return volume.ID, nil
}
