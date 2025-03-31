package titles

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetNewTitles(c *gin.Context) {
	limit := 10

	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "лимит должен быть числом"})
			return
		}
	}

	var titles []struct {
		ID     uint
		Name   string
		Author string
		Team   string
	}

	h.DB.Raw(
		`SELECT t.name, a.name AS author, teams.name AS team
		FROM titles AS t
		INNER JOIN authors AS a ON t.author_id = a.id
		INNER JOIN teams ON t.team_id = teams.id
		ORDER BY t.created_at DESC
		LIMIT ?`, limit,
	).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "новые тайтлы не найдены"}) // мало ли
		return
	}

	c.JSON(200, &titles)
}
