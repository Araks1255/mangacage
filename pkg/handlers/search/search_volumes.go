package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) SearchVolumes(c *gin.Context) {
	query := strings.ToLower(c.Param("query"))

	type result struct {
		Title  string `gorm:"column:name"`
		Volume string `gorm:"column:name"`
	}

	var results []result
	h.DB.Raw(`SELECT titles.name, volumes.name FROM volumes
	INNER JOIN titles ON volumes.title_id = titles.id
	WHERE volumes.name ILIKE ? AND NOT volumes.on_moderation`,
		fmt.Sprintf("%%%s%%", query),
	).Scan(&results)

	if len(results) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено томов по вашему запросу"})
		return
	}

	response := make(map[int]result, len(results))
	for i := 0; i < len(results); i++ {
		response[i] = results[i]
	}

	c.JSON(200, &response)
}
