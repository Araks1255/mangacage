package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchTitles(query string, limit int) (titles *[]models.TitleDTO, err error) {
	var result []models.TitleDTO

	err = h.DB.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			a.name AS author, a.id AS author_id,
			MAX(teams.name) AS team, MAX(teams.id) AS team_id,
			ARRAY_AGG(g.name) AS genres
		FROM
			titles AS t
			INNER JOIN authors AS a ON a.id = t.author_id
			LEFT JOIN teams ON teams.id = t.team_id
			INNER JOIN title_genres AS tg ON t.id = tg.title_id
			INNER JOIN genres AS g ON g.id = tg.genre_id
		WHERE
			lower(t.name) ILIKE lower(?)
		GROUP BY
			t.id, a.id
		LIMIT ?`,
		fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
