package titles

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetMostPopularTitles(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	type title struct {
		Name  string
		Views uint
	}

	var titles []title
	h.DB.Raw(
		`SELECT titles.name, user_viewed_chapters.chapter_id, COUNT(user_viewed_chapters.chapter_id) AS views
		FROM user_viewed_chapters
		INNER JOIN chapters ON user_viewed_chapters.chapter_id = chapters.id
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		GROUP BY titles.name, user_viewed_chapters.chapter_id
		ORDER BY views DESC
		LIMIT ?`, limit).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено самых читаемых тайтлов"})
		return
	}

	response := make(map[int]title, len(titles))
	for i := 0; i < len(titles); i++ {
		response[i] = titles[i]
	}

	c.JSON(200, &response)
}
