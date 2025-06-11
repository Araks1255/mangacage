package search

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	titlesHelpers "github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	"gorm.io/gorm"
)

func SearchTitles(db *gorm.DB, query string, limit int) (titles *[]models.TitleDTO, err error) {
	var result []models.TitleDTO

	err = titlesHelpers.GetTitle(db).
		Where("lower(name) ILIKE lower(?)", query).
		Limit(limit).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
