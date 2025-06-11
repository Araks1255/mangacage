package search

import (
	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

func SearchChapters(db *gorm.DB, query string, limit int) (chapters *[]models.ChapterDTO, err error) {
	var result []models.ChapterDTO

	err = db.Raw(
		`SELECT
			c.id, c.created_at, c.name, c.description, c.number_of_pages,
			v.name AS volume, v.id AS volume_id, t.name AS title, t.id AS title_id
		FROM
			chapters AS c
			INNER JOIN volumes AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			lower(c.name) ILIKE lower(?)
		LIMIT ?`,
		query, limit,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
