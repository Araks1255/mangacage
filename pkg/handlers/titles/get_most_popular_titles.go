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
		Name    string
		TitleID uint `json:"-"`
		Views   uint
	}

	var titles []title
	h.DB.Raw(`SELECT titles.name, user_viewed_titles.title_id, COUNT(user_viewed_titles.title_id) AS views
		FROM user_viewed_titles
		INNER JOIN titles ON user_viewed_titles.title_id = titles.id
		WHERE NOT titles.on_moderation
		GROUP BY titles.name, user_viewed_titles.title_id
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
