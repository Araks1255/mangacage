package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchChapters(query string, limit int) (chapters *[]models.ChapterDTO, quantity int) {
	var result []models.ChapterDTO

	h.DB.Raw(
		`SELECT c.id, c.created_at, c.name, c.description, c.number_of_pages,
		v.name AS volume, v.id AS volume_id, t.name AS title, t.id AS title_id
		FROM chapters AS c
		INNER JOIN volumes AS v ON v.id = c.volume_id
		INNER JOIN titles AS t ON t.id = v.title_id
		WHERE lower(c.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
