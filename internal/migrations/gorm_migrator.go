package migrations

import (
	"log"
	"os"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

func GormMigrate(db *gorm.DB) error {
	result := db.Exec(
		`CREATE TABLE users (
    		id BIGSERIAL PRIMARY KEY,
   			user_name TEXT,
    		team_id BIGINT
		)`,
	)
	if result.Error != nil {
		log.Println(result.Error)
	}

	if result = db.Exec(
		`CREATE TABLE teams (
    		id BIGSERIAL PRIMARY KEY,
    		name TEXT,
    		creator_id BIGINT,
    		moderator_id BIGINT
		)`,
	); result.Error != nil {
		log.Println(result.Error)
	}

	err := db.AutoMigrate(
		&models.Role{},
		&models.Genre{},
	)
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&models.Team{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.Author{},
		&models.Title{},
		&models.Volume{},
		&models.Chapter{},
	); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.TitleOnModeration{},
		&models.VolumeOnModeration{},
		&models.ChapterOnModeration{},
		&models.UserOnModeration{},
		&models.TeamOnModeration{},
	); err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&models.TeamJoinRequest{},
		&models.UserViewedChapter{},
	); err != nil {
		return err
	}

	createGetRecentlyUpdatedTitlesFunctionSQL, err := os.ReadFile("internal/migrations/sql/create_get_recently_updated_titles.sql")
	if err != nil {
		return err
	}

	if result := db.Exec(string(createGetRecentlyUpdatedTitlesFunctionSQL)); result.Error != nil {
		return result.Error
	}

	return nil
}
