package search

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h handler) SearchTitles(c *gin.Context) {
	query := c.Param("query")

	type result struct {
		Title  string `gorm:"column:name"`
		Author string `gorm:"column:name"`
	}

	var results []result
	h.DB.Raw(`SELECT titles.name, authors.name FROM titles
		INNER JOIN authors ON titles.author_id = authors.id
		WHERE lower(titles.name) ILIKE lower(?)
		AND NOT titles.on_moderation`,
		fmt.Sprintf("%%%s%%", query),
	).Scan(&results)

	if len(results) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено тайтлов по вашему запросу"})
		return
	}

	response := make(map[int]result, len(results))
	for i := 0; i < len(results); i++ {
		response[i] = results[i]
	}

	c.JSON(200, &response)
}
