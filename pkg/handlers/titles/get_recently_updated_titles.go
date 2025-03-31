package titles

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetRecentlyUpdatedTitles(c *gin.Context) {
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
		`SELECT t.id, t.name, a.name AS author, teams.name AS team
		FROM titles AS t
		INNER JOIN volumes AS v ON t.id = v.title_id
		INNER JOIN chapters AS c ON v.id = c.volume_id
		INNER JOIN authors AS a ON a.id = t.author_id
		INNER JOIN teams ON teams.id = t.team_id
		ORDER BY c.updated_at DESC
		LIMIT ?`, limit,
	).Scan(&titles)

	if len(titles) == 0 { // Ну мало ли
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено недавно обновлённых тайтлов"})
		return
	}

	c.JSON(200, &titles)
}
