package search

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) Search(c *gin.Context) {
	query := c.Query("query")
	searchingType := c.Query("type")
	limit, err := strconv.ParseUint(c.Query("limit"), 10, 32)

	if searchingType == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "не указана область поиска"})
		return
	}

	if query == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "отсутсвует поисковой запрос"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
		return
	}

	query = fmt.Sprintf("%%%s%%", query)

	var result any

	switch searchingType {
	case "titles":
		result, err = SearchTitles(h.DB, query, int(limit))

	case "chapters":
		result, err = SearchChapters(h.DB, query, int(limit))

	case "volumes":
		result, err = SearchVolumes(h.DB, query, int(limit))

	case "teams":
		result, err = SearchTeams(h.DB, query, int(limit))

	case "authors":
		result, err = SearchAuthors(h.DB, query, int(limit))

	default:
		c.AbortWithStatusJSON(400, gin.H{"error": "недопустимая область поиска"})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
