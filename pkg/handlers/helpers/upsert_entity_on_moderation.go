package helpers

import "gorm.io/gorm"

func UpsertEntityOnModeration[T upsertableEntityOnModeration](db *gorm.DB, entityOnModeration T, entityOnModerationID uint) (err error) {
	if entityOnModerationID == 0 {
		err = db.Create(entityOnModeration).Error
	} else {
		err = db.Model(entityOnModeration).Updates(entityOnModeration).Error
	}
	return err
}
