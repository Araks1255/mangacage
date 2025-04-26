package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchAuthors(query string, limit int) (authors *[]models.AuthorDTO, quantity int) {
	var result []models.AuthorDTO

	h.DB.Raw(
		`SELECT a.id, a.name, a.about, ARRAY_AGG(g.name)::TEXT[] AS genres
		FROM authors AS a
		INNER JOIN author_genres AS ag ON ag.author_id = a.id
		INNER JOIN genres AS g ON g.id = ag.genre_id
		WHERE lower(a.name) ILIKE lower(?)
		GROUP BY a.id
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
