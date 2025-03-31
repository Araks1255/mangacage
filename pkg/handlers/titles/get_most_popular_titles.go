package titles

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetMostPopularTitles(c *gin.Context) {
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
		Views  uint
	}

	h.DB.Raw(
		`SELECT t.id, t.name, a.name AS author, uvs.chapter_id, COUNT(uvs.chapter_id) AS views
		FROM user_viewed_chapters AS uvs 
		INNER JOIN chapters AS c ON uvs.chapter_id = c.id
		INNER JOIN volumes AS v ON c.volume_id = v.id
		INNER JOIN titles AS t ON v.title_id = t.id
		INNER JOIN authors AS a ON a.id = t.author_id
		GROUP BY t.id, t.name, a.name, uvs.chapter_id
		ORDER BY views DESC
		LIMIT ?`, limit).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено самых читаемых тайтлов"})
		return
	}

	c.JSON(200, &titles)
}
