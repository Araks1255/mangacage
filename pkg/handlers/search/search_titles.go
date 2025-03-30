package search

import (
	"github.com/lib/pq"
	"fmt"
	"time"
)

type Title struct {
	ID uint
	CreatedAt time.Time
	Name string
	Description string
	Author string
	Team string
	Genres pq.StringArray
}

func (h handler) SearchTitles(query string, limit int) (titles *[]Title, quantity int) {
	var result []Title

	h.DB.Raw(
		`SELECT t.id, t.created_at, t.name, t.description, authors.name AS author, teams.name AS team,
		(
			SELECT ARRAY(
				SELECT genres.name FROM genres
				INNER JOIN title_genres ON genres.id = title_genres.genre_id
				INNER JOIN titles ON title_genres.title_id = titles.id
				WHERE titles.id = t.id
			) AS genres
		) FROM titles AS t
		INNER JOIN authors ON authors.id = t.author_id
		LEFT JOIN teams ON teams.id = t.team_id
		WHERE lower(t.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
