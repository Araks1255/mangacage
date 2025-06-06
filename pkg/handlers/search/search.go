package search

import (
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

	var result any

	switch searchingType {
	case "titles":
		result, err = h.SearchTitles(query, int(limit))

	case "chapters":
		result, err = h.SearchChapters(query, int(limit))

	case "volumes":
		result, err = h.SearchVolumes(query, int(limit))

	case "teams":
		result, err = h.SearchTeams(query, int(limit))

	case "authors":
		result, err = h.SearchAuthors(query, int(limit))

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
