package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) SearchTeams(c *gin.Context) {
	query := strings.ToLower(c.Param("query"))

	var teams []string
	h.DB.Raw("SELECT name FROM teams WHERE name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Scan(&teams)
	if len(teams) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено команд перевода по вашему запросу"})
		return
	}

	response := make(map[int]string, len(teams))
	for i := 0; i < len(teams); i++ {
		response[i] = teams[i]
	}

	c.JSON(200, &response)
}
