package moderation

import (
	"database/sql"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateVolumeOnModerationOptions struct {
	ExistingID uint
}

func CreateVolumeOnModeration(db *gorm.DB, titleID, userID uint, opts ...CreateVolumeOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	volume := models.VolumeOnModeration{
		Name:      uuid.New().String(),
		TitleID:   sql.NullInt64{Int64: int64(titleID), Valid: true},
		CreatorID: userID,
	}

	if len(opts) != 0 && opts[0].ExistingID != 0 {
		volume.ExistingID = sql.NullInt64{Int64: int64(opts[0].ExistingID), Valid: true}
	}

	if result := db.Create(&volume); result.Error != nil {
		return 0, result.Error
	}

	return volume.ID, nil
}
