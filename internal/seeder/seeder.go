package seeder

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, mongoDB *mongo.Database, mode string) error {
	switch mode {
	case "test":
		if err := seedRoles(db); err != nil {
			return err
		}
		if err := seedGenres(db); err != nil {
			return err
		}
		if err := seedTags(db); err != nil {
			return err
		}
		if err := seedEntities(db, mongoDB); err != nil {
			return err
		}

	case "prod":
		if err := seedRoles(db); err != nil {
			return err
		}

	default:
		return errors.New("Неккоректный тип сида")
	}

	return nil
}
