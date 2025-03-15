package titles

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetNewTitles(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	type title struct {
		Name   string
		Author string
		Team   string
	}

	var titles []title

	h.DB.Raw(
		`SELECT titles.name, authors.name AS author, teams.name AS team FROM titles
		INNER JOIN authors ON titles.author_id = authors.id
		INNER JOIN teams ON titles.team_id = teams.id
		WHERE NOT titles.on_moderation
		ORDER BY titles.created_at DESC
		LIMIT ?`, limit).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error":"новые тайтлы не найдены"}) // мало ли
		return
	}

	response := make(map[int]title)
	for i := 0; i < len(titles); i++ {
		response[i] = titles[i]
	}

	c.JSON(200, response)
}
