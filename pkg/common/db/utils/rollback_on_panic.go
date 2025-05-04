package utils

import "gorm.io/gorm"

func RollbackOnPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}
}
