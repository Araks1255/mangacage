package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) SearchChapters(c *gin.Context) {
	query := strings.ToLower(c.Param("query"))

	type result struct {
		Title   string `gorm:"column:name"`
		Volume  string `gorm:"column:name"`
		Chapter string `gorm:"column:name"`
	}

	var results []result
	h.DB.Raw(`SELECT titles.name, volumes.name, chapters.name FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE chapters.name ILIKE ? AND NOT chapters.on_moderation`,
		fmt.Sprintf("%%%s%%", query),
	).Scan(&results)

	if len(results) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено глав по вашему запросу"})
		return
	}

	response := make(map[int]result, len(results))
	for i := 0; i < len(results); i++ {
		response[i] = results[i]
	}

	c.JSON(200, &response)
}
