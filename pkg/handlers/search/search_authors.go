package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) SearchAuthors(c *gin.Context) {
	query := strings.ToLower(c.Param("query"))

	var authors []string
	h.DB.Raw("SELECT name FROM authors WHERE name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Scan(&authors)
	if len(authors) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено авторов по вашему запросу"})
		return
	}

	response := make(map[int]string, len(authors))
	for i := 0; i < len(authors); i++ {
		response[i] = authors[i]
	}

	c.JSON(200, &response)
}
