package helpers

import (
	"gorm.io/gorm"
)

func UpsertEntityChanges[T upsertableEntityOnModeration](db *gorm.DB, entityOnModeration T, entityOnModerationExistingID uint) (err error) {
	var entityOnModerationID *uint

	err = db.Model(entityOnModeration).
		Select("id").
		Where("existing_id = ?", entityOnModerationExistingID).
		Scan(&entityOnModerationID).Error

	if err != nil {
		return err
	}

	if entityOnModerationID == nil {
		err = db.Create(entityOnModeration).Error
	} else {
		err = db.Model(entityOnModeration).
			Where("id = ?", entityOnModerationID).
			Updates(entityOnModeration).
			Error
		entityOnModeration.SetID(*entityOnModerationID)
	}

	return err
}
