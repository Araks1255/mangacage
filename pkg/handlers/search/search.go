package search

import (
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

	var (
		result any
		quantity int
	)

	switch searchingType {
	case "title":
		result, quantity = h.SearchTitles(query, limit)
	case "chapter":
		result, quantity = h.SearchChapters(query, limit)
	}

	if quantity == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "по вашему запросу ничего не найдено"})
		return
	}

	c.JSON(200, result)
}
