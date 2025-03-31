package search

import (
	"fmt"

	"github.com/lib/pq"
)

type Author struct {
	ID     uint
	Name   string
	About  string
	Genres pq.StringArray `gorm:"type:TEXT[]"`
}

func (h handler) SearchAuthors(query string, limit int) (authors *[]Author, quantity int) {
	var result []Author

	h.DB.Raw(
		`SELECT a.id, a.name, a.about,
		(
			SELECT ARRAY(
				SELECT genres.name FROM genres
				INNER JOIN author_genres ON genres.id = author_genres.genre_id
				WHERE author_genres.author_id = a.id
			) AS genres
		) FROM authors AS a
		WHERE lower(a.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
