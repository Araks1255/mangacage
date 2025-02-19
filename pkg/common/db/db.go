package db

import (
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

	db.AutoMigrate(&models.Title{}, &models.Chapter{}, &models.User{}, &models.Team{}, &models.Genre{}, &models.Author{}, &models.Role{})
	return db, nil
}
