package search

import (
	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

func SearchVolumes(db *gorm.DB, query string, limit int) (volumes *[]models.VolumeDTO, err error) {
	var result []models.VolumeDTO

	err = db.Raw(
		`SELECT
			v.id, v.created_at, v.name, v.description,
			t.name AS title, t.id AS title_id
		FROM
			volumes AS v
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			lower(v.name) ILIKE lower(?)
		LIMIT ?`,
		query, limit,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
