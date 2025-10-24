package migrations

import (
	"context"

	"gorm.io/gorm"
)

func Migrate(ctx context.Context, db *gorm.DB, pathToMediaDir string) error {
	if err := gormMigrate(db); err != nil {
		return err
	}
	if err := fsMigrate(pathToMediaDir); err != nil {
		return err
	}
	return nil
}
