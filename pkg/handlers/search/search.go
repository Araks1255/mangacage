package search

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) Search(c *gin.Context) {
	query := c.Query("query")
	searchingType := c.Query("type")
	limit, err := strconv.Atoi(c.Query("limit"))

	if query == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "отсутсвует поисковой запрос"})
		return
	}

	if err != nil || limit == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "лимит результатов должен быть числом и быть больше нуля"})
		return
	}

	var result any

	switch searchingType {
	case "titles":
		result, err = h.SearchTitles(query, limit)

	case "chapters":
		result, err = h.SearchChapters(query, limit)

	case "volumes":
		result, err = h.SearchVolumes(query, limit)

	case "teams":
		result, err = h.SearchTeams(query, limit)

	case "authors":
		result, err = h.SearchAuthors(query, limit)

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
