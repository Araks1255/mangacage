package search

import (
	"fmt"
	"time"
)

type Volume struct {
	ID          uint
	CreatedAt   time.Time
	Name        string
	Description string
	Title       string
}

func (h handler) SearchVolumes(query string, limit int) (volumes *[]Volume, quantity int) {
	var result []Volume

	h.DB.Raw(
		`SELECT v.id, v.created_at, v.name, v.description, titles.name AS title
		FROM volumes AS v
		INNER JOIN titles ON titles.id = v.title_id
		WHERE lower(v.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
