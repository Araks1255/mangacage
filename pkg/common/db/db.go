package db

import (
	"github.com/Araks1255/mangabrad/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dbUrl string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Title{}, &models.Chapter{}, &models.User{}, &models.Team{}, &models.Genre{}, &models.Author{})
	return db, nil
}
