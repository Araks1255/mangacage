package favorites

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFavoriteTitles(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var titles []struct {
		Name   string
		Author string
	}

	h.DB.Raw(`SELECT titles.name, authors.name AS author FROM titles
		INNER JOIN user_favorite_titles ON titles.id = user_favorite_titles.title_id
		INNER JOIN authors ON authors.id = titles.author_id
		WHERE user_favorite_titles.user_id = ?
		AND NOT titles.on_moderation`, claims.ID).Scan(&titles)

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено избранных тайтлов"})
		return
	}

	c.JSON(200, &titles)
}
