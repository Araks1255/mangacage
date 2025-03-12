package users

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetViewedTitles(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	type title struct {
		Name   string
		Author string
	}

	var titles []title

	h.DB.Raw(`SELECT titles.name, authors.name AS author FROM titles
		INNER JOIN authors ON titles.author_id = authors.id
		INNER JOIN user_viewed_titles ON titles.id = user_viewed_titles.title_id
		INNER JOIN users ON user_viewed_titles.user_id = users.id
		WHERE users.id = ?`, claims.ID).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "нет читаемых тайтлов"})
		return
	}

	response := make(map[int]title, len(titles))
	for i := 0; i < len(titles); i++ {
		response[i] = titles[i]
	}

	c.JSON(200, response)
}
