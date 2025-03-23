package db

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbUrl string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Title{}, &models.Chapter{}, &models.User{}, &models.Team{}, &models.Genre{}, &models.Author{}, &models.Role{},
		&models.TitleOnModeration{}, &models.VolumeOnModeration{}, &models.ChapterOnModeration{}, &models.UserOnModeration{}, &models.TeamOnModeration{},
	)

	if result := db.Exec("INSERT INTO roles (name) VALUES ('user'), ('moder'), ('admin'), ('team_leader'), ('translater')"); result.Error != nil { // Создание необходимых для работы ролей
		log.Println(result.Error)
	}

	return db, nil
}
