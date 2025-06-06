package testhelpers

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateVolumeOptions struct { // Вообше нет смысмла сейчас создавать отдельную структуру, но на будущее не помешает
	ModeratorID uint
}

func CreateVolume(db *gorm.DB, titleID, teamID, creatorID uint, opts ...CreateVolumeOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	volume := models.Volume{
		Name:      uuid.New().String(),
		TitleID:   titleID,
		CreatorID: creatorID,
		TeamID:    teamID,
	}

	if len(opts) != 0 && opts[0].ModeratorID != 0 {
		volume.ModeratorID = &opts[0].ModeratorID
	}

	if result := db.Create(&volume); result.Error != nil {
		return 0, result.Error
	}

	return volume.ID, nil
}
