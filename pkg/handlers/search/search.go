package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h handler) Search(c *gin.Context) {
	table := strings.ToLower(c.Param("type"))
	query := strings.ToLower(c.Param("query"))

	switch table {
	case "titles":
	case "volumes":
		var volumes []string
		h.DB.Raw("SELECT name FROM volumes WHERE name ILIKE ? AND NOT on_moderation", fmt.Sprintf("%%%s%%", query)).Scan(&volumes)
		if len(volumes) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено томов по запросу"})
			return
		}

		var volumesTitles []string
		h.DB.Raw(
			"SELECT titles.name FROM titles INNER JOIN volumes ON titles.id = volumes.title_id WHERE volumes.name = ANY(?) ORDER BY ARRAY_POSITION(ARRAY[?]::text[], volumes.name)",
			pq.Array(volumes),
			pq.Array(volumes),
		).Scan(&volumesTitles)

		result := make(map[string]string, len(volumes))

		for i := 0; i < len(volumes); i++ {
			result[volumes[i]] = volumesTitles[i]
		}

		c.JSON(200, result)

	case "chapters":

	}
}
