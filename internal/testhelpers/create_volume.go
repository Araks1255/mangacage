package testhelpers

import (
	"database/sql"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateVolumeOptions struct { // Вообше нет смысмла сейчас создавать отдельную структуру, но на будущее не помешает
	ModeratorID uint
}

func CreateVolume(db *gorm.DB, titleID, creatorID uint, opts ...CreateVolumeOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	volume := models.Volume{
		Name:      uuid.New().String(),
		TitleID:   titleID,
		CreatorID: creatorID,
	}

	if len(opts) != 0 && opts[0].ModeratorID != 0 {
		volume.ModeratorID = sql.NullInt64{Int64: int64(opts[0].ModeratorID), Valid: true}
	}

	if result := db.Create(&volume); result.Error != nil {
		return 0, result.Error
	}

	return volume.ID, nil
}
