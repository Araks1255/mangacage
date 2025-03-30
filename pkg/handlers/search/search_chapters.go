package search

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Chapter struct {
	ID uint
	CreatedAt time.Time
	Name string
	Description string
	NumberOfPages int
	Volume string
	Title string
}

func (h handler) SearchChapters(query string, limit int) (chapters *[]Chapter, quantity string){
	var result []Chapter

	h.DB.Raw(
		`SELECT c.id, c.created_at, c.name, c.description, c.number_of_pages,
		volumes.name AS volume, titles.name AS title
		FROM chapters AS c
		INNER JOIN volumes ON volumes.id = c.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE lower(c.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}	
