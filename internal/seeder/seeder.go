package seeder

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, mongoDB *mongo.Database, mode string) error {
	tx := db.Begin()
	defer tx.Rollback()

	switch mode {
	case "test":
		if err := seedRoles(tx); err != nil {
			return err
		}
		if err := seedGenres(tx); err != nil {
			return err
		}
		if err := seedTags(tx); err != nil {
			return err
		}

	case "prod":
		if err := seedRoles(tx); err != nil {
			return err
		}

	default:
		return errors.New("Неккоректный тип сида")
	}

	tx.Commit()

	return nil
}
