package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Migrate(ctx context.Context, db *gorm.DB, mongoDB *mongo.Database) error {
	if err := gormMigrate(db); err != nil {
		return err
	}
	if err := mongoMigrate(ctx, mongoDB); err != nil {
		return err
	}
	return nil
}
