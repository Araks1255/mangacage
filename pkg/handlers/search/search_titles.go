package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchTitles(query string, limit int) (titles *[]models.TitleDTO, quantity int) {
	var result []models.TitleDTO

	h.DB.Raw(
		`SELECT t.id, t.created_at, t.name, t.description,
		a.name AS author, a.id AS author_id, teams.name AS team, teams.id AS team_id,
		(
			SELECT ARRAY(
				SELECT g.name FROM genres AS g
				INNER JOIN title_genres AS tg ON g.id = tg.genre_id
				WHERE tg.title_id = t.id
			) AS genres
		) FROM titles AS t
		INNER JOIN authors AS a ON a.id = t.author_id
		LEFT JOIN teams ON teams.id = t.team_id
		WHERE lower(t.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
