package search

import (
	"fmt"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) Search(c *gin.Context) {
	table := strings.ToLower(c.Param("type"))
	query := strings.ToLower(c.Param("query"))

	if table != "titles" && table != "chapters" && table != "genres" && table != "authors" && table != "teams" {
		c.AbortWithStatusJSON(404, gin.H{"error":"недопустимая область поиска"})
		return
	}

	var results []string
	h.DB.Raw(fmt.Sprintf("SELECT name FROM %s WHERE name ILIKE ?", table), fmt.Sprintf("%%%s%%", query)).Scan(&results)

	response := utils.ConvertToMap(results)
	c.JSON(200, response)
}